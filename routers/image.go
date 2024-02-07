package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/naponte/Udemy_Go_React_MongoDB/awsgo"
	"github.com/naponte/Udemy_Go_React_MongoDB/bd"
	"github.com/naponte/Udemy_Go_React_MongoDB/models"
)

type readSeeker struct {
	io.Reader
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func UploadImage(ctx context.Context, uploadType string, req events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.Status = 400

	userID := claim.ID.Hex()

	var filename string
	var user models.User

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))

	switch uploadType {
	case "A":
		filename = "avatars/" + userID + ".jpg"
		user.Avatar = filename
	case "B":
		filename = "banners/" + userID + ".jpg"
		user.Banner = filename
	}

	mediaType, params, err := mime.ParseMediaType(req.Headers["Content-Type"])
	if err != nil {
		response.Status = 500
		response.Message = err.Error()
		return response
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		response.Message = "Contant-Type header must bo 'multipart/' type"
		return response
	}

	body, err := base64.StdEncoding.DecodeString(req.Body)
	if err != nil {
		response.Message = "Error decoding body " + err.Error()
		response.Status = 500
		return response
	}

	mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	p, err := mr.NextPart()
	if err != nil && err != io.EOF {
		response.Message = err.Error()
		response.Status = 500
		return response
	}
	if err != io.EOF {
		if p.FileName() != "" {
			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, p); err != nil {
				response.Message = err.Error()
				response.Status = 500
				return response
			}

			session, err := session.NewSession(&aws.Config{
				Region: aws.String("us-east-1"),
			})
			if err != nil {
				response.Message = err.Error()
				response.Status = 500
				return response
			}

			uploader := s3manager.NewUploader(session)
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: bucket,
				Key:    aws.String(filename),
				Body:   &readSeeker{buf},
			})
			if err != nil {
				response.Message = err.Error()
				response.Status = 500
				return response
			}
		}
	}

	status, err := bd.UpdateProfile(user, userID)
	if err != nil || !status {
		response.Status = 400
		response.Message = "Error updating user profile " + err.Error()
		return response
	}

	response.Status = 200
	response.Message = "Image upload OK!"
	return response
}

func GetImage(ctx context.Context, uploadType string, req events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.Status = 400

	ID := req.QueryStringParameters["id"]
	if len(ID) < 1 {
		response.Message = "ID is required"
		return response
	}

	profile, err := bd.SearchProfile(ID)
	if err != nil {
		response.Message = "User " + ID + " not found " + err.Error()
		return response
	}

	var fileName string
	switch uploadType {
	case "A":
		fileName = profile.Avatar
	case "B":
		fileName = profile.Banner
	}
	fmt.Println("fileName: " + fileName)

	svc := s3.NewFromConfig(awsgo.Cfg)
	file, err := downloadFromS3(ctx, svc, fileName)
	if err != nil {
		response.Status = 500
		response.Message = "Error downloading file " + err.Error()
		return response
	}

	response.CustomResp = &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       file.String(),
		Headers: map[string]string{
			"Content-Type":        "application/octet-stream",
			"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", fileName),
		},
	}

	return response
}

func downloadFromS3(ctx context.Context, svc *s3.Client, fileName string) (*bytes.Buffer, error) {
	bucket := ctx.Value(models.Key("bucketName")).(string)
	obj, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}

	defer obj.Body.Close()

	fmt.Println("bucketName: " + bucket)

	file, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(file)
	return buffer, nil
}
