package settings

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/mcuadros/go-version"
	"go.uber.org/zap"
)

var (
	settingsInstance *AppSettings
)

const (
	SETTINGS_FILENAME         = "settings.json"
	TITLE_JSON_FILENAME       = "titles.json"
	VERSIONS_JSON_FILENAME    = "versions.json"
	SLM_VERSION               = "1.9.0"
	DEFAULT_TITLES_JSON_URL   = "https://tinfoil.media/repo/db/titles.json"
	DEFAULT_VERSIONS_JSON_URL = "https://raw.githubusercontent.com/blawar/titledb/master/versions.json"
	SLM_VERSION_URL           = "https://raw.githubusercontent.com/firebat20/switch-library-manager/master/version.json"
)

const (
	TEMPLATE_TITLE_ID    = "TITLE_ID"
	TEMPLATE_TITLE_NAME  = "TITLE_NAME"
	TEMPLATE_DLC_NAME    = "DLC_NAME"
	TEMPLATE_VERSION     = "VERSION"
	TEMPLATE_REGION      = "REGION"
	TEMPLATE_VERSION_TXT = "VERSION_TXT"
	TEMPLATE_TYPE        = "TYPE"
)

type OrganizeOptions struct {
	CreateFolderPerGame        bool   `json:"create_folder_per_game"`
	DlcFolder                  string `json:"dlc_folder"`
	UpdatesFolder              string `json:"updates_folder"`
	RenameFiles                bool   `json:"rename_files"`
	DeleteEmptyFolders         bool   `json:"delete_empty_folders"`
	DeleteOldUpdateFiles       bool   `json:"delete_old_update_files"`
	FolderNameTemplate         string `json:"folder_name_template"`
	SwitchSafeFileNames        bool   `json:"switch_safe_file_names"`
	FileNameTemplate           string `json:"file_name_template"`
	ProcessWhenMissingBaseGame bool   `json:"process_when_missing_base_game"`
}

type AppSettings struct {
	VersionsJsonUrl        string          `json:"versions_json_url"`
	VersionsEtag           string          `json:"versions_etag"`
	TitlesJsonUrl          string          `json:"titles_json_url"`
	TitlesEtag             string          `json:"titles_etag"`
	Prodkeys               string          `json:"prod_keys"`
	Folder                 string          `json:"folder"`
	ScanFolders            []string        `json:"scan_folders"`
	GUI                    bool            `json:"gui"`
	Debug                  bool            `json:"debug"`
	CheckForMissingUpdates bool            `json:"check_for_missing_updates"`
	CheckForMissingDLC     bool            `json:"check_for_missing_dlc"`
	HideMissingGames       bool            `json:"hide_missing_games"`
	HideDemoGames          bool            `json:"hide_demo_games"`
	OrganizeOptions        OrganizeOptions `json:"organize_options"`
	ScanRecursively        bool            `json:"scan_recursively"`
	GuiPagingSize          int             `json:"gui_page_size"`
	IgnoreDLCUpdates       bool            `json:"ignore_dlc_updates"`
	IgnoreDLCTitleIds      []string        `json:"ignore_dlc_title_ids"`
	IgnoreUpdateTitleIds   []string        `json:"ignore_update_title_ids"`
	IgnoreFileTypes        []string        `json:"ignore_file_types"`
}

func ReadSettingsAsJSON(baseFolder string) string {
	if _, err := os.Stat(filepath.Join(baseFolder, SETTINGS_FILENAME)); err != nil {
		saveDefaultSettings(baseFolder)
	}
	file, _ := os.Open(filepath.Join(baseFolder, SETTINGS_FILENAME))
	bytes, _ := io.ReadAll(file)
	return string(bytes)
}

func ReadSettings(baseFolder string) *AppSettings {
	if settingsInstance != nil {
		return settingsInstance
	}
	settingsInstance = &AppSettings{Debug: false, GuiPagingSize: 100, ScanFolders: []string{},
		OrganizeOptions: OrganizeOptions{SwitchSafeFileNames: true}, Prodkeys: "", IgnoreDLCTitleIds: []string{"01007F600B135007"}}
	if _, err := os.Stat(filepath.Join(baseFolder, SETTINGS_FILENAME)); err == nil {
		file, err := os.Open(filepath.Join(baseFolder, SETTINGS_FILENAME))
		if err != nil {
			zap.S().Warnf("Missing or corrupted config file, creating a new one")
			return saveDefaultSettings(baseFolder)
		} else {
			_ = json.NewDecoder(file).Decode(&settingsInstance)
			settingsInstance = verifySettings(baseFolder, settingsInstance)
			return settingsInstance
		}
	} else {
		return saveDefaultSettings(baseFolder)
	}
}

func verifySettings(baseFolder string, settings *AppSettings) *AppSettings {
	// check so titles json url is set, if not revert to default
	if settings.TitlesJsonUrl == "" {
		settings.TitlesJsonUrl = DEFAULT_TITLES_JSON_URL
	}
	// check so version json url is set, if not revert to default
	if settings.VersionsJsonUrl == "" {
		settings.VersionsJsonUrl = DEFAULT_VERSIONS_JSON_URL
	}

	// check to the title json file exists, if it does not, revert ETAG
	if _, err := os.Stat(filepath.Join(baseFolder, TITLE_JSON_FILENAME)); err != nil {
		settings.TitlesEtag = "W/\"a5b02845cf6bd61:0\""
	}
	// check to the version json file exists, if it does not, revert ETAG
	if _, err := os.Stat(filepath.Join(baseFolder, VERSIONS_JSON_FILENAME)); err != nil {
		settings.VersionsEtag = "W/\"2ef50d1cb6bd61:0\""
	}

	return settings
}

func saveDefaultSettings(baseFolder string) *AppSettings {
	settingsInstance = &AppSettings{
		TitlesJsonUrl:          DEFAULT_TITLES_JSON_URL,
		TitlesEtag:             "W/\"a5b02845cf6bd61:0\"",
		VersionsJsonUrl:        DEFAULT_VERSIONS_JSON_URL,
		VersionsEtag:           "W/\"2ef50d1cb6bd61:0\"",
		Folder:                 "",
		Prodkeys:               "",
		ScanFolders:            []string{},
		IgnoreUpdateTitleIds:   []string{},
		IgnoreDLCTitleIds:      []string{},
		IgnoreDLCUpdates:       false,
		IgnoreFileTypes:        []string{},
		GUI:                    true,
		GuiPagingSize:          100,
		CheckForMissingUpdates: true,
		CheckForMissingDLC:     true,
		HideMissingGames:       false,
		HideDemoGames:          false,
		ScanRecursively:        true,
		Debug:                  false,
		OrganizeOptions: OrganizeOptions{
			RenameFiles:         false,
			CreateFolderPerGame: false,
			DlcFolder:           "",
			UpdatesFolder:       "",
			FolderNameTemplate:  fmt.Sprintf("{%v}", TEMPLATE_TITLE_NAME),
			FileNameTemplate: fmt.Sprintf("{%v} ({%v})[{%v}][v{%v}]", TEMPLATE_TITLE_NAME, TEMPLATE_DLC_NAME,
				TEMPLATE_TITLE_ID, TEMPLATE_VERSION),
			DeleteEmptyFolders:         false,
			SwitchSafeFileNames:        true,
			DeleteOldUpdateFiles:       false,
			ProcessWhenMissingBaseGame: false,
		},
	}
	return SaveSettings(settingsInstance, baseFolder)
}

func SaveSettings(settings *AppSettings, baseFolder string) *AppSettings {
	file, _ := json.MarshalIndent(settings, "", " ")
	_ = os.WriteFile(filepath.Join(baseFolder, SETTINGS_FILENAME), file, 0644)
	settingsInstance = settings
	return settings
}

func CheckForUpdates() (bool, error) {

	localVer := SLM_VERSION

	res, err := http.Get(SLM_VERSION_URL)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	remoteValues := map[string]string{}
	err = json.Unmarshal(body, &remoteValues)
	if err != nil {
		return false, err
	}

	remoteVer := remoteValues["version"]

	if version.CompareSimple(remoteVer, localVer) > 0 {
		return true, nil
	}

	return false, nil
}
