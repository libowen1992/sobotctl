package minioOps

import (
	"context"
	"fmt"
	miniot "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/spf13/cobra"
	"log"
	"sobotctl/global"
)

func NewMinioOpsCmd() *cobra.Command {
	minioOpsCmd := &cobra.Command{
		Use:   "minio",
		Short: "minio manager",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("minio")
			endpoint := "62.234.40.29:19000"
			accessKeyID := "storage"
			secretAccessKey := "WmmZ6OiqjMioxueMsJyShofGafT4TAXW"
			useSSL := false

			// Initialize minio client object.
			minioClient, err := miniot.New(endpoint, &miniot.Options{
				Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
				Secure: useSSL,
			})
			if err != nil {
				global.Logger.Error("err create minio client", err)
				return
			}

			global.Logger.Info(minioClient)
			log.Printf("%#v\n", minioClient) // minioClient is now set up

			ctx, cancel := context.WithCancel(context.Background())

			defer cancel()

			objectCh := minioClient.ListObjects(ctx, "storage", miniot.ListObjectsOptions{
				Recursive: true,
			})
			for object := range objectCh {
				if object.Err != nil {
					global.Logger.Error(object.Err)
					return
				}
				fmt.Println(object)
			}
		},
	}
	return minioOpsCmd
}
