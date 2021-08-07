package bridge

import (
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/object"
	"os"
)

func uploadPhotoToVk(vkToken string, vkClub int, filename string) (object.PhotosPhoto, error) {
	vk := api.NewVK(vkToken)
	file, err := os.Open(filename)
	if err != nil {
		return object.PhotosPhoto{}, err
	}
	photos, err := vk.UploadGroupWallPhoto(vkClub, file)
	if err != nil {
		return object.PhotosPhoto{}, err
	}
	return photos[0], err
}
