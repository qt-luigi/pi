package pi

import (
	"fmt"
	"io/ioutil"
	"testing"
)

var graphTests = []struct {
	name     string
	input    []string
	exitCode int
}{
	{
		name:     "create graph - not specify id",
		input:    []string{"graphs", "create", "--name", "test-name", "--type", "int", "--unit", "commits", "--color", "shibafu", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "create graph - not specify name",
		input:    []string{"graphs", "create", "--graph-id", "test-id", "--type", "int", "--unit", "commits", "--color", "shibafu", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "create graph - not specify type",
		input:    []string{"graphs", "create", "--graph-id", "test-id", "--name", "test-name", "--unit", "commits", "--color", "shibafu", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "create graph - not specify unit",
		input:    []string{"graphs", "create", "--graph-id", "test-id", "--name", "test-name", "--type", "int", "--color", "shibafu", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "create graph - not specify color",
		input:    []string{"graphs", "create", "--graph-id", "test-id", "--name", "test-name", "--type", "int", "--unit", "commits", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "create graph - not specify username",
		input:    []string{"graphs", "create", "--graph-id", "test-id", "--name", "test-name", "--type", "int", "--unit", "commits", "--color", "shibafu"},
		exitCode: 1,
	},
	{
		name:     "create graph - color is invalid",
		input:    []string{"graphs", "create", "--graph-id", "test-id", "--name", "test-name", "--type", "int", "--unit", "commits", "--color", "rainbow", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "create graph - type is invalid",
		input:    []string{"graphs", "create", "--graph-id", "test-id", "--name", "test-name", "--type", "string", "--unit", "commits", "--color", "shibafu", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "create graph - self-sufficient is invalid",
		input:    []string{"graphs", "create", "--graph-id", "test-id", "--name", "test-name", "--type", "int", "--unit", "commits", "--color", "shibafu", "--username", "c-know", "--self-sufficient", "yes"},
		exitCode: 1,
	},
	{
		name:     "get graph definition - not specify username",
		input:    []string{"graphs", "get"},
		exitCode: 1,
	},
	{
		name:     "get svg graph url - not specify username",
		input:    []string{"graphs", "svg", "--graph-id", "test-id"},
		exitCode: 1,
	},
	{
		name:     "get svg graph url - not specify id",
		input:    []string{"graphs", "svg", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "update graph - not specify id",
		input:    []string{"graphs", "update", "--name", "test-name", "--unit", "commits", "--color", "shibafu", "--username", "c-know", "--purge-cache-urls", "http://example.com/a", "--purge-cache-urls", "http://example.com/b"},
		exitCode: 1,
	},
	{
		name:     "update graph - not specify username",
		input:    []string{"graphs", "update", "--name", "test-name", "--unit", "commits", "--color", "shibafu", "--graph-id", "test-id", "--purge-cache-urls", "http://example.com/a", "--purge-cache-urls", "http://example.com/b"},
		exitCode: 1,
	},
	{
		name:     "update graph - invalid color name",
		input:    []string{"graphs", "update", "--name", "test-name", "--unit", "commits", "--color", "rainbow", "--username", "c-know", "--graph-id", "test-id", "--purge-cache-urls", "http://example.com/a", "--purge-cache-urls", "http://example.com/b"},
		exitCode: 1,
	},
	{
		name:     "update graph - invalid self-sufficient",
		input:    []string{"graphs", "update", "--name", "test-name", "--unit", "commits", "--color", "shibafu", "--username", "c-know", "--graph-id", "test-id", "--self-sufficient", "ok", "--purge-cache-urls", "http://example.com/a", "--purge-cache-urls", "http://example.com/b"},
		exitCode: 1,
	},
	{
		name:     "update graph - purge cache urls limit over",
		input:    []string{"graphs", "update", "--name", "test-name", "--unit", "commits", "--color", "shibafu", "--username", "c-know", "--graph-id", "test-id", "--purge-cache-urls", "http://example.com/a", "--purge-cache-urls", "http://example.com/b", "--purge-cache-urls", "http://example.com/c", "--purge-cache-urls", "http://example.com/d", "--purge-cache-urls", "http://example.com/e", "--purge-cache-urls", "http://example.com/f"},
		exitCode: 1,
	},
	{
		name:     "get graph detail url - not specify username",
		input:    []string{"graphs", "detail", "--graph-id", "test-id"},
		exitCode: 1,
	},
	{
		name:     "get graph detail url - not specify id",
		input:    []string{"graphs", "detail", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "get graph list page url - not specify username",
		input:    []string{"graphs", "list"},
		exitCode: 1,
	},
	{
		name:     "delete graph - not specify username",
		input:    []string{"graphs", "delete", "--graph-id", "test-id"},
		exitCode: 1,
	},
	{
		name:     "delete graph - not specify id",
		input:    []string{"graphs", "delete", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "get graph pixels - not specify username",
		input:    []string{"graphs", "pixels", "--graph-id", "test-id"},
		exitCode: 1,
	},
	{
		name:     "get graph pixels - not specify id",
		input:    []string{"graphs", "pixels", "--username", "c-know"},
		exitCode: 1,
	},
	{
		name:     "get graph stats - not specify username",
		input:    []string{"graphs", "stats", "--graph-id", "test-id"},
		exitCode: 1,
	},
	{
		name:     "get graph stats - not specify id",
		input:    []string{"graphs", "stats", "--username", "c-know"},
		exitCode: 1,
	},
}

func TestGraph(t *testing.T) {
	for _, tt := range graphTests {
		exitCode := (&CLI{
			ErrStream: ioutil.Discard,
			OutStream: ioutil.Discard,
		}).Run(tt.input)
		if exitCode != tt.exitCode {
			t.Errorf("%s(exitCode): out=%d want=%d", tt.name, exitCode, tt.exitCode)
		}
	}
}

func TestGenerateCreateGraphRequest(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testType := "int"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testIsSecret := true
	testPublishOptionalData := true
	cmd := &createGraphCommand{
		Username:            testUsername,
		ID:                  testID,
		Name:                testName,
		Unit:                testUnit,
		Type:                testType,
		Color:               testColor,
		Timezone:            testTimezone,
		SelfSufficient:      testSelfSufficient,
		Secret:              &testIsSecret,
		PublishOptionalData: &testPublishOptionalData,
	}

	// run
	req, err := generateCreateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "POST" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs", afterAPIBaseEnv, testUsername) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request body. %s", err)
	}
	if string(b) != fmt.Sprintf(`{"id":"%s","name":"%s","unit":"%s","type":"%s","color":"%s","timezone":"%s","selfSufficient":"%s","isSecret":%t,"publishOptionalData":%t}`, testID, testName, testUnit, testType, testColor, testTimezone, testSelfSufficient, testIsSecret, testPublishOptionalData) {
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateCreateGraphRequestWithoutIsSecret(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testType := "int"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	cmd := &createGraphCommand{
		Username:       testUsername,
		ID:             testID,
		Name:           testName,
		Unit:           testUnit,
		Type:           testType,
		Color:          testColor,
		Timezone:       testTimezone,
		SelfSufficient: testSelfSufficient,
	}

	// run
	req, err := generateCreateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "POST" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs", afterAPIBaseEnv, testUsername) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request body. %s", err)
	}
	if string(b) != fmt.Sprintf(`{"id":"%s","name":"%s","unit":"%s","type":"%s","color":"%s","timezone":"%s","selfSufficient":"%s"}`, testID, testName, testUnit, testType, testColor, testTimezone, testSelfSufficient) {
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateGetGraphRequest(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	cmd := &getGraphsCommand{
		Username: testUsername,
	}

	// run
	req, err := generateGetGraphsRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "GET" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs", afterAPIBaseEnv, testUsername) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			t.Errorf("Failed to read request body. %s", err)
		}
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateSVGUrlNoParam(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testDate := ""
	testMode := ""
	cmd := &graphSVGCommand{
		Username: testUsername,
		ID:       testID,
		Date:     testDate,
		Mode:     testMode,
	}

	// run
	url, err := generateSVGUrl(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if url != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected url. %s", url)
	}
}

func TestGenerateSVGUrlDateSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testDate := "20190101"
	testMode := ""
	testAppearance := ""
	cmd := &graphSVGCommand{
		Username:   testUsername,
		ID:         testID,
		Date:       testDate,
		Mode:       testMode,
		Appearance: testAppearance,
	}

	// run
	url, err := generateSVGUrl(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if url != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s?date=%s", afterAPIBaseEnv, testUsername, testID, testDate) {
		t.Errorf("Unexpected url. %s", url)
	}
}

func TestGenerateSVGUrlModeSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testDate := ""
	testMode := "short"
	testAppearance := ""
	cmd := &graphSVGCommand{
		Username:   testUsername,
		ID:         testID,
		Date:       testDate,
		Mode:       testMode,
		Appearance: testAppearance,
	}

	// run
	url, err := generateSVGUrl(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if url != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s?mode=%s", afterAPIBaseEnv, testUsername, testID, testMode) {
		t.Errorf("Unexpected url. %s", url)
	}
}

func TestGenerateSVGUrlAppearanceSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testDate := ""
	testMode := ""
	testAppearance := "dark"
	cmd := &graphSVGCommand{
		Username:   testUsername,
		ID:         testID,
		Date:       testDate,
		Mode:       testMode,
		Appearance: testAppearance,
	}

	// run
	url, err := generateSVGUrl(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if url != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s?appearance=%s", afterAPIBaseEnv, testUsername, testID, testAppearance) {
		t.Errorf("Unexpected url. %s", url)
	}
}

func TestGenerateSVGUrlDateAndModeParamSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testDate := "20190101"
	testMode := "short"
	testAppearance := ""
	cmd := &graphSVGCommand{
		Username:   testUsername,
		ID:         testID,
		Date:       testDate,
		Mode:       testMode,
		Appearance: testAppearance,
	}

	// run
	url, err := generateSVGUrl(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if url != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s?date=%s&mode=%s", afterAPIBaseEnv, testUsername, testID, testDate, testMode) {
		t.Errorf("Unexpected url. %s", url)
	}
}

func TestGenerateSVGUrlModeAndAppearanceParamSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testDate := ""
	testMode := "short"
	testAppearance := "dark"
	cmd := &graphSVGCommand{
		Username:   testUsername,
		ID:         testID,
		Date:       testDate,
		Mode:       testMode,
		Appearance: testAppearance,
	}

	// run
	url, err := generateSVGUrl(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if url != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s?mode=%s&appearance=%s", afterAPIBaseEnv, testUsername, testID, testMode, testAppearance) {
		t.Errorf("Unexpected url. %s", url)
	}
}

func TestGenerateSVGUrlDateAndAppearanceParamSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testDate := "20190331"
	testMode := ""
	testAppearance := "dark"
	cmd := &graphSVGCommand{
		Username:   testUsername,
		ID:         testID,
		Date:       testDate,
		Mode:       testMode,
		Appearance: testAppearance,
	}

	// run
	url, err := generateSVGUrl(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if url != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s?date=%s&appearance=%s", afterAPIBaseEnv, testUsername, testID, testDate, testAppearance) {
		t.Errorf("Unexpected url. %s", url)
	}
}

func TestGenerateSVGUrlAllParamSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testDate := "20190331"
	testMode := "short"
	testAppearance := "dark"
	cmd := &graphSVGCommand{
		Username:   testUsername,
		ID:         testID,
		Date:       testDate,
		Mode:       testMode,
		Appearance: testAppearance,
	}

	// run
	url, err := generateSVGUrl(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if url != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s?date=%s&mode=%s&appearance=%s", afterAPIBaseEnv, testUsername, testID, testDate, testMode, testAppearance) {
		t.Errorf("Unexpected url. %s", url)
	}
}

func TestGenerateUpdateGraphRequestWithSecretAndPublish(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, _, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/"}
	testIsSecret := true
	testIsPublish := true
	cmd := &updateGraphCommand{
		Username:       testUsername,
		ID:             testID,
		Name:           testName,
		Unit:           testUnit,
		Color:          testColor,
		Timezone:       testTimezone,
		PurgeCacheURLs: testPurgeCacheURLs,
		SelfSufficient: testSelfSufficient,
		Secret:         &testIsSecret,
		Publish:        &testIsPublish,
	}

	// run
	_, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err == nil {
		t.Errorf("Error should have occurs.")
	}
}

func TestGenerateUpdateGraphRequestWithSecret(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/"}
	testIsSecret := true
	cmd := &updateGraphCommand{
		Username:       testUsername,
		ID:             testID,
		Name:           testName,
		Unit:           testUnit,
		Color:          testColor,
		Timezone:       testTimezone,
		PurgeCacheURLs: testPurgeCacheURLs,
		SelfSufficient: testSelfSufficient,
		Secret:         &testIsSecret,
	}

	// run
	req, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "PUT" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request body. %s", err)
	}
	if string(b) != fmt.Sprintf(`{"name":"%s","unit":"%s","color":"%s","timezone":"%s","purgeCacheURLs":["%s"],"selfSufficient":"%s","isSecret":%t}`, testName, testUnit, testColor, testTimezone, testPurgeCacheURLs[0], testSelfSufficient, testIsSecret) {
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateUpdateGraphRequestWithSomeParams(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/"}
	testIsSecret := true
	cmd := &updateGraphCommand{
		Username: testUsername,
		ID:       testID,
		Name:     testName,
		// Unit not spesify
		// Unit:           testUnit,
		Color:          testColor,
		Timezone:       testTimezone,
		PurgeCacheURLs: testPurgeCacheURLs,
		SelfSufficient: testSelfSufficient,
		Secret:         &testIsSecret,
	}

	// run
	req, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "PUT" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request body. %s", err)
	}
	if string(b) != fmt.Sprintf(`{"name":"%s","color":"%s","timezone":"%s","purgeCacheURLs":["%s"],"selfSufficient":"%s","isSecret":%t}`, testName, testColor, testTimezone, testPurgeCacheURLs[0], testSelfSufficient, testIsSecret) {
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateUpdateGraphRequestWithPublish(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/"}
	testIsPublish := true
	cmd := &updateGraphCommand{
		Username:       testUsername,
		ID:             testID,
		Name:           testName,
		Unit:           testUnit,
		Color:          testColor,
		Timezone:       testTimezone,
		PurgeCacheURLs: testPurgeCacheURLs,
		SelfSufficient: testSelfSufficient,
		Secret:         nil,
		Publish:        &testIsPublish,
	}

	// run
	req, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "PUT" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request body. %s", err)
	}
	if string(b) != fmt.Sprintf(`{"name":"%s","unit":"%s","color":"%s","timezone":"%s","purgeCacheURLs":["%s"],"selfSufficient":"%s","isSecret":%t}`, testName, testUnit, testColor, testTimezone, testPurgeCacheURLs[0], testSelfSufficient, !testIsPublish) {
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateUpdateGraphRequestWithoutSecretOrPublish(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/"}
	cmd := &updateGraphCommand{
		Username:       testUsername,
		ID:             testID,
		Name:           testName,
		Unit:           testUnit,
		Color:          testColor,
		Timezone:       testTimezone,
		PurgeCacheURLs: testPurgeCacheURLs,
		SelfSufficient: testSelfSufficient,
	}

	// run
	req, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "PUT" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request body. %s", err)
	}
	if string(b) != fmt.Sprintf(`{"name":"%s","unit":"%s","color":"%s","timezone":"%s","purgeCacheURLs":["%s"],"selfSufficient":"%s"}`, testName, testUnit, testColor, testTimezone, testPurgeCacheURLs[0], testSelfSufficient) {
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateUpdateGraphRequestWithHideAndPublishOptionalData(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, _, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/"}
	testIsSecret := true
	testIsPublish := true
	cmd := &updateGraphCommand{
		Username:            testUsername,
		ID:                  testID,
		Name:                testName,
		Unit:                testUnit,
		Color:               testColor,
		Timezone:            testTimezone,
		PurgeCacheURLs:      testPurgeCacheURLs,
		SelfSufficient:      testSelfSufficient,
		PublishOptionalData: &testIsSecret,
		HideOptionalData:    &testIsPublish,
	}

	// run
	_, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err == nil {
		t.Errorf("Error should have occurs.")
	}
}

func TestGenerateUpdateGraphRequestWithPublishOptionalData(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/"}
	testPublishOptionalData := true
	cmd := &updateGraphCommand{
		Username:            testUsername,
		ID:                  testID,
		Name:                testName,
		Unit:                testUnit,
		Color:               testColor,
		Timezone:            testTimezone,
		PurgeCacheURLs:      testPurgeCacheURLs,
		SelfSufficient:      testSelfSufficient,
		PublishOptionalData: &testPublishOptionalData,
	}

	// run
	req, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "PUT" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request body. %s", err)
	}
	if string(b) != fmt.Sprintf(`{"name":"%s","unit":"%s","color":"%s","timezone":"%s","purgeCacheURLs":["%s"],"selfSufficient":"%s","publishOptionalData":%t}`, testName, testUnit, testColor, testTimezone, testPurgeCacheURLs[0], testSelfSufficient, testPublishOptionalData) {
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateUpdateGraphRequestWithHideOptionalData(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/"}
	testHideOptionalData := true
	cmd := &updateGraphCommand{
		Username:         testUsername,
		ID:               testID,
		Name:             testName,
		Unit:             testUnit,
		Color:            testColor,
		Timezone:         testTimezone,
		PurgeCacheURLs:   testPurgeCacheURLs,
		SelfSufficient:   testSelfSufficient,
		HideOptionalData: &testHideOptionalData,
	}

	// run
	req, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "PUT" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request body. %s", err)
	}
	if string(b) != fmt.Sprintf(`{"name":"%s","unit":"%s","color":"%s","timezone":"%s","purgeCacheURLs":["%s"],"selfSufficient":"%s","publishOptionalData":%t}`, testName, testUnit, testColor, testTimezone, testPurgeCacheURLs[0], testSelfSufficient, !testHideOptionalData) {
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateUpdateGraphRequestWithoutHideOrPublishOptionalData(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/"}
	cmd := &updateGraphCommand{
		Username:       testUsername,
		ID:             testID,
		Name:           testName,
		Unit:           testUnit,
		Color:          testColor,
		Timezone:       testTimezone,
		PurgeCacheURLs: testPurgeCacheURLs,
		SelfSufficient: testSelfSufficient,
	}

	// run
	req, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "PUT" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	b, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request body. %s", err)
	}
	if string(b) != fmt.Sprintf(`{"name":"%s","unit":"%s","color":"%s","timezone":"%s","purgeCacheURLs":["%s"],"selfSufficient":"%s"}`, testName, testUnit, testColor, testTimezone, testPurgeCacheURLs[0], testSelfSufficient) {
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateUpdateGraphRequestOver5URLsSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, _, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testName := "test-name"
	testUnit := "commits"
	testColor := "ajisai"
	testTimezone := "asia/Tokyo"
	testSelfSufficient := "none"
	testPurgeCacheURLs := []string{"http://test.example.com/1", "http://test.example.com/2", "http://test.example.com/3", "http://test.example.com/4", "http://test.example.com/5", "http://test.example.com/6"}
	cmd := &updateGraphCommand{
		Username:       testUsername,
		ID:             testID,
		Name:           testName,
		Unit:           testUnit,
		Color:          testColor,
		Timezone:       testTimezone,
		PurgeCacheURLs: testPurgeCacheURLs,
		SelfSufficient: testSelfSufficient,
	}

	// run
	_, err := generateUpdateGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err == nil {
		t.Errorf("Error should have occurs.")
	}
}

func TestGenerateDeleteGraphRequest(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	cmd := &deleteGraphCommand{
		Username: testUsername,
		ID:       testID,
	}

	// run
	req, err := generateDeleteGraphRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "DELETE" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			t.Errorf("Failed to read request body. %s", err)
		}
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateGetGraphPixelsRequest(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	cmd := &getGraphPixelsCommand{
		Username: testUsername,
		ID:       testID,
	}

	// run
	req, err := generateGetGraphPixelsRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "GET" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s/pixels", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			t.Errorf("Failed to read request body. %s", err)
		}
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateGetGraphPixelsRequestFromParamSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testFrom := "20190101"
	cmd := &getGraphPixelsCommand{
		Username: testUsername,
		ID:       testID,
		From:     testFrom,
	}

	// run
	req, err := generateGetGraphPixelsRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "GET" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s/pixels?from=%s", afterAPIBaseEnv, testUsername, testID, testFrom) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			t.Errorf("Failed to read request body. %s", err)
		}
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateGetGraphPixelsRequestToParamSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testFrom := ""
	testTo := "20191231"
	cmd := &getGraphPixelsCommand{
		Username: testUsername,
		ID:       testID,
		From:     testFrom,
		To:       testTo,
	}

	// run
	req, err := generateGetGraphPixelsRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "GET" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s/pixels?to=%s", afterAPIBaseEnv, testUsername, testID, testTo) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			t.Errorf("Failed to read request body. %s", err)
		}
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateGetGraphPixelsRequestBothParamSpecified(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	testFrom := "20190101"
	testTo := "20191231"
	cmd := &getGraphPixelsCommand{
		Username: testUsername,
		ID:       testID,
		From:     testFrom,
		To:       testTo,
	}

	// run
	req, err := generateGetGraphPixelsRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "GET" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s/pixels?from=%s&to=%s", afterAPIBaseEnv, testUsername, testID, testFrom, testTo) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			t.Errorf("Failed to read request body. %s", err)
		}
		t.Errorf("Unexpected request body. %s", string(b))
	}
}

func TestGenerateGetGraphStatsRequest(t *testing.T) {
	// prepare
	beforeAPIBaseEnv, beforeTokenEnv, afterAPIBaseEnv, _ := prepare()

	testUsername := "c-know"
	testID := "test-id"
	cmd := &getGraphStatsCommand{
		Username: testUsername,
		ID:       testID,
	}

	// run
	req, err := generateGetGraphStatsRequest(cmd)

	// cleanup
	cleanup(beforeAPIBaseEnv, beforeTokenEnv)

	// assertion
	if err != nil {
		t.Errorf("Unexpected error occurs. %s", err)
	}
	if req.Method != "GET" {
		t.Errorf("Unexpected request method. %s", req.Method)
	}
	if req.URL.String() != fmt.Sprintf("https://%s/v1/users/%s/graphs/%s/stats", afterAPIBaseEnv, testUsername, testID) {
		t.Errorf("Unexpected request path. %s", req.URL.String())
	}
	if req.Body != nil {
		b, err := ioutil.ReadAll(req.Body)
		defer req.Body.Close()
		if err != nil {
			t.Errorf("Failed to read request body. %s", err)
		}
		t.Errorf("Unexpected request body. %s", string(b))
	}
}
