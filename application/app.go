package application

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"restaurant/common"
	"restaurant/domain"
	"restaurant/entity"
	"restaurant/pkgs"

	officedocCommon "github.com/unidoc/unioffice/common"
	"github.com/unidoc/unioffice/common/license"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
)

func init() {
	err := license.SetMeteredKey("416172dd38b123b5d3efbbf78d22b2014ee572e958a6bdb9b5f3c2864dbb1925")
	if err != nil {
		log.Printf("ERROR to get license: %v", err)
		panic(err)
	}

}

type AppService struct {
	db      *domain.AppRepo
	account AccountService
}

func NewAppService() *AppService {
	return &AppService{
		db:      domain.NewAppRepo(),
		account: *NewAccountService(),
	}
}

func (as *AppService) GetApp(ctx context.Context, id int64) (*entity.AppResp, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, nil
	}

	param, err := as.application2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (as *AppService) SearchById(ctx context.Context, id int64) (map[string]interface{}, error) {
	apps, total, err := as.db.SearchById(ctx, id)
	if err != nil {
		return nil, err
	}

	params := make([]entity.AppResp, 0)
	for _, app := range apps {
		param := &entity.AppResp{
			Id:          app.Id,
			UserId:      app.UserId,
			Name:        app.Name,
			Address:     app.Address,
			Describe:    app.Describe,
			Phone:       app.Phone,
			PostCode:    app.PostCode,
			PostName:    app.PostName,
			WorkBeginAt: app.WorkBeginAt,
			WorkEndAt:   app.WorkEndAt,
			HaveVege:    app.HaveVege,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt,
			UpdatedAt:   app.UpdatedAt,
		}

		bodyPics := make([]string, 0)
		var appPics []string
		err = json.Unmarshal([]byte(app.Pictures), &appPics)
		if err != nil {
			return nil, err
		}

		for _, pic := range appPics {
			signedUrl, _ := pkgs.SignedUrl(pic, false)
			bodyPics = append(bodyPics, signedUrl)
		}

		log.Println(bodyPics)

		param.PicsUrl = bodyPics
		param.Pictures = appPics

		params = append(params, *param)
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (as *AppService) GetApps(ctx context.Context, ids []int64) ([]*entity.Application, error) {
	app, err := as.db.GetApps(ctx, ids)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (as *AppService) GetAppsOrder(ctx context.Context, ids []int64) ([]*entity.AppResp, error) {
	apps, err := as.db.GetAppsOrder(ctx, ids)
	if err != nil {
		return nil, err
	}

	log.Println(apps)

	params := make([]*entity.AppResp, 0)
	for _, app := range apps {
		param := &entity.AppResp{
			Id:          app.Id,
			UserId:      app.UserId,
			Name:        app.Name,
			Address:     app.Address,
			Describe:    app.Describe,
			Phone:       app.Phone,
			PostCode:    app.PostCode,
			PostName:    app.PostName,
			WorkBeginAt: app.WorkBeginAt,
			WorkEndAt:   app.WorkEndAt,
			HaveVege:    app.HaveVege,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt,
			UpdatedAt:   app.UpdatedAt,
		}

		bodyPics := make([]string, 0)
		var appPics []string
		err = json.Unmarshal([]byte(app.Pictures), &appPics)
		if err != nil {
			return nil, err
		}

		for _, pic := range appPics {
			signedUrl, _ := pkgs.SignedUrl(pic, false)
			bodyPics = append(bodyPics, signedUrl)
		}

		log.Println(bodyPics)

		param.PicsUrl = bodyPics
		param.Pictures = appPics

		params = append(params, param)
	}

	return params, nil
}

func (a *AppService) application2Resp(app *entity.Application) (*entity.AppResp, error) {
	param := &entity.AppResp{
		Id:          app.Id,
		UserId:      app.UserId,
		Name:        app.Name,
		Address:     app.Address,
		Describe:    app.Describe,
		Phone:       app.Phone,
		PostCode:    app.PostCode,
		PostName:    app.PostName,
		WorkBeginAt: app.WorkBeginAt,
		WorkEndAt:   app.WorkEndAt,
		HaveVege:    app.HaveVege,
		Status:      app.Status,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}

	bodyPics := make([]string, 0)
	var appPics []string
	err := json.Unmarshal([]byte(app.Pictures), &appPics)
	if err != nil {
		return nil, err
	}

	for _, pic := range appPics {
		signedUrl, _ := pkgs.SignedUrl(pic, false)
		bodyPics = append(bodyPics, signedUrl)
	}

	log.Println(bodyPics)
	param.PicsUrl = bodyPics
	param.Pictures = appPics

	return param, nil
}

func (a *AppService) appResp2Application(app *entity.AppResp) (*entity.Application, error) {
	param := &entity.Application{
		Id:          app.Id,
		UserId:      app.UserId,
		Name:        app.Name,
		Address:     app.Address,
		Describe:    app.Describe,
		Phone:       app.Phone,
		PostCode:    app.PostCode,
		PostName:    app.PostName,
		WorkBeginAt: app.WorkBeginAt,
		WorkEndAt:   app.WorkEndAt,
		HaveVege:    app.HaveVege,
		Status:      app.Status,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}

	pics, err := json.Marshal(app.Pictures)
	if err != nil {
		return nil, err
	}

	param.Pictures = string(pics)

	return param, nil
}

// Create 创建
func (a *AppService) Create(ctx context.Context, param entity.AppResp, userId int64, status entity.Status) (*entity.AppResp, error) {
	app, err := a.GetAppByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	currentTime := common.NowInLocal()

	pics := param.Pictures
	data, err := json.Marshal(pics)
	if err != nil {
		return nil, err
	}

	params := entity.Application{
		Pictures:    string(data),
		Name:        param.Name,
		Address:     param.Address,
		Describe:    param.Describe,
		Phone:       param.Phone,
		PostCode:    param.PostCode,
		PostName:    param.PostName,
		WorkBeginAt: param.WorkBeginAt,
		WorkEndAt:   param.WorkEndAt,
		HaveVege:    param.HaveVege,
	}
	if app != nil {
		params.Id = app.Id
		params.UserId = app.UserId
		params.Status = status
		params.UpdatedAt = currentTime
		params.CreatedAt = app.CreatedAt
		tempApp, err := a.Update(ctx, userId, app.Id, &params)
		if err != nil {
			return nil, err
		}

		ret, err := a.application2Resp(tempApp)
		if err != nil {
			return nil, err
		}
		return ret, nil
	}

	params.UserId = userId
	params.CreatedAt = currentTime
	params.UpdatedAt = currentTime
	params.Status = status

	tempApp, err := a.db.Create(ctx, &params)
	if err != nil {
		return nil, err
	}

	ret, err := a.application2Resp(tempApp)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

// List 列表
func (a *AppService) List(ctx context.Context, offset, limit int) (map[string]interface{}, error) {
	apps, total, err := a.db.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	params := make([]entity.AppResp, 0)
	for _, app := range apps {
		param := &entity.AppResp{
			Id:          app.Id,
			UserId:      app.UserId,
			Name:        app.Name,
			Address:     app.Address,
			Describe:    app.Describe,
			Phone:       app.Phone,
			PostCode:    app.PostCode,
			PostName:    app.PostName,
			WorkBeginAt: app.WorkBeginAt,
			WorkEndAt:   app.WorkEndAt,
			HaveVege:    app.HaveVege,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt,
			UpdatedAt:   app.UpdatedAt,
		}

		bodyPics := make([]string, 0)
		var appPics []string
		err = json.Unmarshal([]byte(app.Pictures), &appPics)
		if err != nil {
			return nil, err
		}

		for _, pic := range appPics {
			signedUrl, _ := pkgs.SignedUrl(pic, false)
			bodyPics = append(bodyPics, signedUrl)
		}

		log.Println(bodyPics)

		param.PicsUrl = bodyPics
		param.Pictures = appPics

		params = append(params, *param)
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

// List 列表
func (a *AppService) SearchList(ctx context.Context, offset, limit int) (map[string]interface{}, error) {
	apps, total, err := a.db.SearchList(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	params := make([]entity.AppResp, 0)
	for _, app := range apps {
		param := &entity.AppResp{
			Id:          app.Id,
			UserId:      app.UserId,
			Name:        app.Name,
			Address:     app.Address,
			Describe:    app.Describe,
			Phone:       app.Phone,
			PostCode:    app.PostCode,
			PostName:    app.PostName,
			WorkBeginAt: app.WorkBeginAt,
			WorkEndAt:   app.WorkEndAt,
			HaveVege:    app.HaveVege,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt,
			UpdatedAt:   app.UpdatedAt,
		}

		bodyPics := make([]string, 0)
		var appPics []string
		err = json.Unmarshal([]byte(app.Pictures), &appPics)
		if err != nil {
			return nil, err
		}

		for _, pic := range appPics {
			signedUrl, _ := pkgs.SignedUrl(pic, false)
			bodyPics = append(bodyPics, signedUrl)
		}

		log.Println(bodyPics)

		param.PicsUrl = bodyPics
		param.Pictures = appPics

		params = append(params, *param)
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (as *AppService) GetAppByUserId(ctx context.Context, userId int64) (*entity.AppResp, error) {
	app, err := as.db.GetByUserId(ctx, userId)
	if err != nil {
		log.Printf("GetAppByUserId_err: %+v", err)
		return nil, err
	}

	if app == nil {
		return nil, nil
	}

	param, err := as.application2Resp(app)
	if err != nil {
		return nil, err
	}
	return param, nil
}

func (as *AppService) ReadApp(ctx context.Context, id int64) (*entity.AppResp, error) {
	app, err := as.db.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if app.Status != entity.StatusAccepted {
		return nil, pkgs.ServerError
	}

	param, err := as.application2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

// Get 获取自己的应用
func (a *AppService) GetSelfApp(ctx context.Context, userId, id int64) (*entity.AppResp, error) {
	app, err := a.db.Get(ctx, id)
	if err != nil {
		log.Printf("Queryapp_Get_err: %+v", err)
		return nil, fmt.Errorf("app not found")
	}
	if app == nil {
		return nil, pkgs.AppNotFound
	}

	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return nil, err
	}

	if app.UserId != userId && user.Status != entity.AdministratorStatus {
		return nil, pkgs.NoPermission
	}

	param, err := a.application2Resp(app)
	if err != nil {
		return nil, err
	}

	return param, nil
}

func (as *AppService) GetAppName(ctx context.Context, userId, id int64) (string, error) {
	app, err := as.GetApp(ctx, id)
	if err != nil {
		return "", err
	}

	if app.UserId != userId {
		return "", pkgs.NoPermission
	}

	return app.Name, nil
}

func (as *AppService) GetAppDescribe(ctx context.Context, userId, id int64) (string, error) {
	app, err := as.GetApp(ctx, id)
	if err != nil {
		return "", err
	}

	if app.UserId != userId {
		return "", pkgs.NoPermission
	}

	return app.Describe, nil
}

func (as *AppService) GetAppAddress(ctx context.Context, userId, id int64) (string, error) {
	app, err := as.GetApp(ctx, id)
	if err != nil {
		return "", err
	}

	if app.UserId != userId {
		return "", pkgs.NoPermission
	}

	return app.Address, nil
}
func (as *AppService) GetAppPhone(ctx context.Context, userId, id int64) (string, error) {
	app, err := as.GetApp(ctx, id)
	if err != nil {
		return "", err
	}

	if app.UserId != userId {
		return "", pkgs.NoPermission
	}

	return app.Phone, nil
}

func (a *AppService) Delete(ctx context.Context, userId, id int64) error {
	app, err := a.db.Get(ctx, id)
	if err != nil {
		return err
	}

	if app.UserId != userId {
		return pkgs.NoPermission
	}

	return a.db.Delete(ctx, id)
}

func (a *AppService) Update(ctx context.Context, userId, id int64, param *entity.Application) (*entity.Application, error) {
	app, err := a.GetApp(ctx, id)
	if err != nil {
		return nil, err
	}

	if app.UserId != userId {
		return nil, pkgs.NoPermission
	}
	return a.db.Update(ctx, param)
}

func (a *AppService) UpdateStatus(ctx context.Context, userId, id int64, status string) error {
	app, err := a.GetApp(ctx, id)
	if err != nil {
		return err
	}

	user, err := a.account.GetAccount(ctx, userId)
	if err != nil {
		return err
	}
	accountStatus := user.Status
	if accountStatus != entity.AdministratorStatus {
		return pkgs.NoPermission
	}

	if status == "通过" {
		app.Status = entity.StatusAccepted
	} else if status == "拒绝" {
		app.Status = entity.StatusRefused
	}

	param, err := a.appResp2Application(app)
	if err != nil {
		return err
	}
	_, err = a.db.Update(ctx, param)

	return err
}

func (as *AppService) SearchByName(ctx context.Context, name string) (map[string]interface{}, error) {
	apps, total, err := as.db.SearchByName(ctx, name)
	if err != nil {
		return nil, err
	}

	params := make([]entity.AppResp, 0)
	for _, app := range apps {
		param := &entity.AppResp{
			Id:          app.Id,
			UserId:      app.UserId,
			Name:        app.Name,
			Address:     app.Address,
			Describe:    app.Describe,
			Phone:       app.Phone,
			PostCode:    app.PostCode,
			PostName:    app.PostName,
			WorkBeginAt: app.WorkBeginAt,
			WorkEndAt:   app.WorkEndAt,
			HaveVege:    app.HaveVege,
			Status:      app.Status,
			CreatedAt:   app.CreatedAt,
			UpdatedAt:   app.UpdatedAt,
		}

		bodyPics := make([]string, 0)
		var appPics []string
		err = json.Unmarshal([]byte(app.Pictures), &appPics)
		if err != nil {
			return nil, err
		}

		for _, pic := range appPics {
			signedUrl, _ := pkgs.SignedUrl(pic, false)
			bodyPics = append(bodyPics, signedUrl)
		}

		log.Println(bodyPics)

		param.PicsUrl = bodyPics
		param.Pictures = appPics

		params = append(params, *param)
	}

	return map[string]interface{}{
		"total": total,
		"apps":  params,
	}, nil
}

func (as *AppService) DownloadDoc(ctx context.Context, userId, id int64) (string, error) {
	user, err := as.account.GetAccount(ctx, userId)
	if err != nil {
		return "", err
	}

	if user.Status != entity.AdministratorStatus {
		return "", pkgs.NoPermission
	}

	app, err := as.db.Get(ctx, id)
	if err != nil {
		return "", err
	}

	params, err := as.application2Resp(app)
	if err != nil {
		return "", fmt.Errorf("download_doc_application_to_response_err: %v", err)
	}

	doc := document.New()
	defer doc.Close()

	tempDir := "../tempImages"
	os.MkdirAll(tempDir, os.ModePerm)
	defer os.RemoveAll(tempDir)

	imgs := []officedocCommon.ImageRef{}
	for i, url := range params.PicsUrl {
		response, err := http.Get(url)
		if err != nil {
			log.Fatalf("Error downloading image: %v", err)
		}
		defer response.Body.Close()

		imgPath := filepath.Join(tempDir, filepath.Base(params.Pictures[i]))
		file, err := os.Create(imgPath)
		if err != nil {
			log.Fatalf("Error creating image file: %v", err)
		}
		defer file.Close()

		_, err = io.Copy(file, response.Body)
		if err != nil {
			log.Fatalf("Error saving image file: %v", err)
		}

		img, err := officedocCommon.ImageFromFile(imgPath)
		if err != nil {
			return "", fmt.Errorf("unable to create image: %s", err)
		}

		imgref, err := doc.AddImage(img)
		if err != nil {
			return "", err
		}
		imgs = append(imgs, imgref)
	}

	para := doc.AddParagraph()
	run := para.AddRun()
	run.AddText(fmt.Sprintf("餐厅名称:\t%s", params.Name))
	run.AddBreak()

	run.AddText("图片:")
	run.AddBreak()
	for _, img := range imgs {
		inl, err := run.AddDrawingInline(img)
		if err != nil {
			return "", fmt.Errorf("unable to add inline image: %s", err)
		}
		inl.SetSize(2*measurement.Inch, 2*measurement.Inch)
		run.AddTab()
	}
	run.AddBreak()

	run.AddText(fmt.Sprintf("详细地址:\t%s", params.Address))
	run.AddBreak()

	run.AddText(fmt.Sprintf("餐厅联系方式:\t%s", params.Phone))
	run.AddBreak()

	run.AddText(fmt.Sprintf("营业时间:\t%s~%s", params.WorkBeginAt, params.WorkEndAt))
	run.AddBreak()

	run.AddText(fmt.Sprintln("是否有素食:\t"))
	if params.HaveVege == 1 {
		run.AddText("是")
	} else {
		run.AddText("否")
	}

	saveFilePath := fmt.Sprintf("%s.docx", params.Name)
	err = doc.SaveToFile(saveFilePath)
	defer os.Remove(saveFilePath)
	if err != nil {
		return "", err
	}

	err = pkgs.UploadOSSByFilePath(saveFilePath, saveFilePath, false)
	if err != nil {
		return "", fmt.Errorf("Upload_OSS_File_err:%v", err)
	}

	signedUrl, _ := pkgs.SignedUrl(saveFilePath, false)
	return signedUrl, nil
}
