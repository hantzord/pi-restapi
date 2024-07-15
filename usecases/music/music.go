package music

import (
	"capstone/constants"
	"capstone/entities"
	musicEntities "capstone/entities/music"
	"capstone/utilities"
	"mime/multipart"
)

type MusicUseCase struct {
	musicInterface musicEntities.RepositoryInterface
}

func NewMusicUseCase(musicInterface musicEntities.RepositoryInterface) *MusicUseCase {
	return &MusicUseCase{
		musicInterface: musicInterface,
	}
}

func (musicUseCase *MusicUseCase) GetAllMusics(metadata entities.Metadata, userId int, search string) ([]musicEntities.Music, error) {
	musics, err := musicUseCase.musicInterface.GetAllMusics(metadata, userId, search)
	if err != nil {
		return []musicEntities.Music{}, err
	}
	return musics, nil
}

func (musicUseCase *MusicUseCase) GetAllMusicsByDoctorId(metadata entities.MetadataFull, userId int) ([]musicEntities.Music, error) {
	musics, err := musicUseCase.musicInterface.GetAllMusicsByDoctorId(metadata, userId)
	if err != nil {
		return []musicEntities.Music{}, err
	}
	return musics, nil
}

func (musicUseCase *MusicUseCase) GetMusicById(musicId int, userId int) (musicEntities.Music, error) {
	music, err := musicUseCase.musicInterface.GetMusicById(musicId, userId)
	if err != nil {
		return musicEntities.Music{}, err
	}
	return music, nil
}

func (musicUseCase *MusicUseCase) GetLikedMusics(metadata entities.Metadata, userId int) ([]musicEntities.Music, error) {
	musics, err := musicUseCase.musicInterface.GetLikedMusics(metadata, userId)
	if err != nil {
		return []musicEntities.Music{}, err
	}
	return musics, nil
}

func (musicUseCase *MusicUseCase) LikeMusic(musicId int, userId int) error {
	err := musicUseCase.musicInterface.LikeMusic(musicId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (musicUseCase *MusicUseCase) UnlikeMusic(musicId int, userId int) error {
	err := musicUseCase.musicInterface.UnlikeMusic(musicId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (musicUseCase *MusicUseCase) CountMusicByDoctorId(doctorId int) (int, error) {
	count, err := musicUseCase.musicInterface.CountMusicByDoctorId(doctorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (musicUseCase *MusicUseCase) CountMusicLikesByDoctorId(doctorId int) (int, error) {
	count, err := musicUseCase.musicInterface.CountMusicLikesByDoctorId(doctorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (musicUseCase *MusicUseCase) CountMusicViewCountByDoctorId(doctorId int) (int, error) {
	count, err := musicUseCase.musicInterface.CountMusicViewCountByDoctorId(doctorId)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (musicUseCase *MusicUseCase) CountMusicViewByMonth(doctorId int, startMonth string, endMonth string) (map[int]int, error) {
	if startMonth == "" || endMonth == "" {
		return map[int]int{}, constants.ErrEmptyInputViewByMonth
	}
	
	count, err := musicUseCase.musicInterface.CountMusicViewByMonth(doctorId, startMonth, endMonth)
	if err != nil {
		return map[int]int{}, err
	}
	return count, nil
}

func (musicUseCase *MusicUseCase) PostMusic(music musicEntities.Music, fileImage *multipart.FileHeader, fileMusic *multipart.FileHeader) (musicEntities.Music, error) {
	if music.Title == "" || music.Singer == "" {
		return musicEntities.Music{}, constants.ErrEmptyInputMusic
	}

	if fileImage != nil {
		secureUrl, err := utilities.UploadImage(fileImage)
		if err != nil {
			return musicEntities.Music{}, constants.ErrUploadImage
		}
		music.ImageUrl = secureUrl
	}

	if fileMusic != nil {
		secureUrl, err := utilities.UploadImage(fileMusic)
		if err != nil {
			return musicEntities.Music{}, constants.ErrUploadImage
		}
		music.MusicUrl = secureUrl
	}
	
	music, err := musicUseCase.musicInterface.PostMusic(music)
	if err != nil {
		return musicEntities.Music{}, err
	}
	return music, nil
}

func (musicUseCase *MusicUseCase) GetMusicByIdForDoctor(musicId int) (musicEntities.Music, error) {
	music, err := musicUseCase.musicInterface.GetMusicByIdForDoctor(musicId)
	if err != nil {
		return musicEntities.Music{}, err
	}
	return music, nil
}

func (musicUseCase *MusicUseCase) EditMusic(music musicEntities.Music, fileImage *multipart.FileHeader) (musicEntities.Music, error) {
	if music.Title == "" || music.Singer == "" {
		return musicEntities.Music{}, constants.ErrEmptyInputMusic
	}

	if fileImage != nil {
		secureUrl, err := utilities.UploadImage(fileImage)
		if err != nil {
			return musicEntities.Music{}, constants.ErrUploadImage
		}
		music.ImageUrl = secureUrl
	}else{
		music.ImageUrl = ""
	}
	
	music, err := musicUseCase.musicInterface.EditMusic(music)
	if err != nil {
		return musicEntities.Music{}, err
	}
	return music, nil
}

func (musicUseCase *MusicUseCase) DeleteMusic(musicId int) error {
	err := musicUseCase.musicInterface.DeleteMusic(musicId)
	if err != nil {
		return err
	}
	return nil
}