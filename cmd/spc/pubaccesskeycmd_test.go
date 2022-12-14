package main

import (
	"fmt"
	"strings"
	"testing"

	"gitlab.com/NebulousLabs/errors"

	"github.com/EvilRedHorse/pubaccess-node/node/api/client"
	"github.com/EvilRedHorse/pubaccess-node/pubaccesskey"
)

const testSkykeyString string = "pubaccesskey:Aa71WcCoKFwVGAVotJh3USAslb8dotVJp2VZRRSAG2QhYRbuTbQhjDolIJ1nOlQ-rWYK29_1xee5?name=test_key1"

// TestSkykeyCommands tests the basic functionality of the spc pubaccesskey commands
// interface. More detailed testing of the pubaccesskey manager is done in the pubaccesskey
// package.
func TestSkykeyCommands(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	groupDir := siacTestDir(t.Name())

	// Define subtests
	subTests := []subTest{
		{name: "TestDuplicateSkykeyAdd", test: testDuplicateSkykeyAdd},
		{name: "TestChangeKeyEntropyKeepName", test: testChangeKeyEntropyKeepName},
		{name: "TestAddKeyTwice", test: testAddKeyTwice},
		{name: "TestDelete", test: testDeleteKey},
		{name: "TestInvalidPubaccesskeyType", test: testInvalidPubaccesskeyType},
		{name: "TestSkykeyGet", test: testSkykeyGet},
		{name: "TestSkykeyGetUsingNameAndID", test: testSkykeyGetUsingNameAndID},
		{name: "TestSkykeyGetUsingNoNameAndNoID", test: testSkykeyGetUsingNoNameAndNoID},
		{name: "TestSkykeyListKeys", test: testSkykeyListKeys},
		{name: "TestSkykeyListKeysDoesntShowPrivateKeys", test: testSkykeyListKeysDoesntShowPrivateKeys},
		{name: "TestSkykeyListKeysAdditionalKeys", test: testSkykeyListKeysAdditionalKeys},
		{name: "TestSkykeyListKeysAdditionalKeysDoesntShowPrivateKeys", test: testSkykeyListKeysAdditionalKeysDoesntShowPrivateKeys},
	}

	// Run tests
	if err := runSubTests(t, groupDir, subTests); err != nil {
		t.Fatal(err)
	}
}

// testDuplicateSkykeyAdd tests that adding with duplicate Pubaccesskey will return
// duplicate name error.
func testDuplicateSkykeyAdd(t *testing.T, c client.Client) {
	err := skykeyAdd(c, testSkykeyString)
	if err != nil {
		t.Fatal(err)
	}

	err = skykeyAdd(c, testSkykeyString)
	if !strings.Contains(err.Error(), pubaccesskey.ErrSkykeyWithIDAlreadyExists.Error()) {
		t.Fatal("Expected duplicate name error but got", err)
	}
}

// testChangeKeyEntropyKeepName tests that adding with changed entropy but the
// same Pubaccesskey will return duplicate name error.
func testChangeKeyEntropyKeepName(t *testing.T, c client.Client) {
	// Change the key entropy, but keep the same name.
	var sk pubaccesskey.Pubaccesskey
	err := sk.FromString(testSkykeyString)
	if err != nil {
		t.Fatal(err)
	}
	sk.Entropy[0] ^= 1 // flip the first byte.

	skString, err := sk.ToString()
	if err != nil {
		t.Fatal(err)
	}

	// This should return a duplicate name error.
	err = skykeyAdd(c, skString)
	if !strings.Contains(err.Error(), pubaccesskey.ErrSkykeyWithNameAlreadyExists.Error()) {
		t.Fatal("Expected duplicate name error but got", err)
	}
}

// testAddKeyTwice tests that creating a Pubaccesskey with the same key name twice
// returns duplicate name error.
func testAddKeyTwice(t *testing.T, c client.Client) {
	// Check that adding same key twice returns an error.
	keyName := "createkey1"
	_, err := skykeyCreate(c, keyName, pubaccesskey.TypePublicID.ToString())
	if err != nil {
		t.Fatal(err)
	}
	_, err = skykeyCreate(c, keyName, pubaccesskey.TypePublicID.ToString())
	if !strings.Contains(err.Error(), pubaccesskey.ErrSkykeyWithNameAlreadyExists.Error()) {
		t.Fatal("Expected error when creating key with same name")
	}
}

// testDeleteKey tests deleting a key.
func testDeleteKey(t *testing.T, c client.Client) {
	// Create a key.
	keyName := "keyToDeleteByName"
	_, err := skykeyCreate(c, keyName, pubaccesskey.TypePublicID.ToString())
	if err != nil {
		t.Fatal(err)
	}

	// Get the Key
	sk, err := c.SkykeyGetByName(keyName)
	if err != nil {
		t.Fatal(err)
	}
	if sk.Name != keyName {
		t.Fatalf("Expected Skykey name %v but got %v", keyName, sk.Name)
	}

	// Verify incorrect param usage will return an error and will not delete the
	// key
	err = skykeyDelete(c, "name", "id")
	if !errors.Contains(err, errBothNameAndIDUsed) {
		t.Fatalf("Unexpected Error: got `%v`, expected `%v`", err, errBothNameAndIDUsed)
	}
	sk, err = c.SkykeyGetByName(keyName)
	if err != nil {
		t.Fatal(err)
	}
	if sk.Name != keyName {
		t.Fatalf("Expected Skykey name %v but got %v", keyName, sk.Name)
	}
	err = skykeyDelete(c, "", "")
	if !errors.Contains(err, errNeitherNameNorIDUsed) {
		t.Fatalf("Unexpected Error: got `%v`, expected `%v`", err, errNeitherNameNorIDUsed)
	}
	sk, err = c.SkykeyGetByName(keyName)
	if err != nil {
		t.Fatal(err)
	}
	if sk.Name != keyName {
		t.Fatalf("Expected Skykey name %v but got %v", keyName, sk.Name)
	}

	// Delete key by name
	err = skykeyDelete(c, keyName, "")
	if err != nil {
		t.Fatal(err)
	}

	// Try and get the key again
	_, err = c.SkykeyGetByName(keyName)
	if err == nil || !strings.Contains(err.Error(), pubaccesskey.ErrNoSkykeysWithThatName.Error()) {
		t.Fatalf("Expected Error to contain %v and got %v", pubaccesskey.ErrNoSkykeysWithThatName, err)
	}

	// Create key again
	_, err = skykeyCreate(c, keyName, pubaccesskey.TypePublicID.ToString())
	if err != nil {
		t.Fatal(err)
	}

	// Get ID
	sk, err = c.SkykeyGetByName(keyName)
	if err != nil {
		t.Fatal(err)
	}

	// Delete key by ID
	err = skykeyDelete(c, "", sk.ID().ToString())
	if err != nil {
		t.Fatal(err)
	}

	// Try and get the key again
	_, err = c.SkykeyGetByName(keyName)
	if err == nil || !strings.Contains(err.Error(), pubaccesskey.ErrNoSkykeysWithThatName.Error()) {
		t.Fatalf("Expected Error to contain %v and got %v", pubaccesskey.ErrNoSkykeysWithThatName, err)
	}
}

// testInvalidPubaccesskeyType tests that invalid cipher types are caught.
func testInvalidPubaccesskeyType(t *testing.T, c client.Client) {
	// Verify invalid type returns error
	_, err := skykeyCreate(c, "createkey2", pubaccesskey.TypeInvalid.ToString())
	if !strings.Contains(err.Error(), pubaccesskey.ErrInvalidPubaccesskeyType.Error()) {
		t.Fatal("Expected error when creating key with invalid type", err)
	}

	// Submitting a blank pubaccesskey type should succeed and default to private
	keyName := "blankType"
	_, err = skykeyCreate(c, keyName, "")
	if err != nil {
		t.Fatal(err)
	}
	sk, err := c.SkykeyGetByName(keyName)
	if err != nil {
		t.Fatal(err)
	}
	if sk.Type != pubaccesskey.TypePrivateID {
		t.Fatal("Skykey type expected to be private")
	}
	// Delete Key to not impact future sub tests
	err = skykeyDelete(c, keyName, "")
	if err != nil {
		t.Fatal(err)
	}

	// Verify a random type gives an error
	_, err = skykeyCreate(c, "", "not a type")
	if !strings.Contains(err.Error(), pubaccesskey.ErrInvalidPubaccesskeyType.Error()) {
		t.Fatal("Expected error when creating key with invalid skykeytpe", err)
	}
}

// testSkykeyGet tests skykeyGet with known key should not return any errors.
func testSkykeyGet(t *testing.T, c client.Client) {
	keyName := "createkey testSkykeyGet"
	newSkykey, err := skykeyCreate(c, keyName, pubaccesskey.TypePublicID.ToString())
	if err != nil {
		t.Fatal(err)
	}

	// Test skykeyGet
	// known key should have no errors.
	getKeyStr, err := skykeyGet(c, keyName, "")
	if err != nil {
		t.Fatal(err)
	}
	if getKeyStr != newSkykey {
		t.Fatal("Expected keys to match")
	}
}

// testSkykeyGetUsingNameAndID tests using both name and id params should return an
// error.
func testSkykeyGetUsingNameAndID(t *testing.T, c client.Client) {
	_, err := skykeyGet(c, "name", "id")
	if err == nil {
		t.Fatal("Expected error when using both name and id")
	}
}

// testSkykeyGetUsingNoNameAndNoID test using neither name or id param should return an
// error.
func testSkykeyGetUsingNoNameAndNoID(t *testing.T, c client.Client) {
	_, err := skykeyGet(c, "", "")
	if err == nil {
		t.Fatal("Expected error when using neither name or id params")
	}
}

// testSkykeyListKeys tests that skykeyListKeys shows key names, ids and keys
func testSkykeyListKeys(t *testing.T, c client.Client) {
	nKeys := 3
	nExtraLines := 3
	keyStrings := make([]string, nKeys)
	keyNames := make([]string, nKeys)
	keyIDs := make([]string, nKeys)

	initSkykeyData(t, c, keyStrings, keyNames, keyIDs)

	keyListString, err := skykeyListKeys(c, true)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < nKeys; i++ {
		if !strings.Contains(keyListString, keyNames[i]) {
			t.Log(keyListString)
			t.Fatal("Missing key name!", i)
		}
		if !strings.Contains(keyListString, keyIDs[i]) {
			t.Log(keyListString)
			t.Fatal("Missing id!", i)
		}
		if !strings.Contains(keyListString, keyStrings[i]) {
			t.Log(keyListString)
			t.Fatal("Missing key!", i)
		}
	}
	keyList := strings.Split(keyListString, "\n")
	if len(keyList) != nKeys+nExtraLines {
		t.Log(keyListString)
		t.Fatalf("Unexpected number of lines/keys %d, Expected %d", len(keyList), nKeys+nExtraLines)
	}
}

// testKskykeyListKeysDoesntShowPrivateKeys tests that skykeyListKeys shows key names, ids and
// doesn't show private keys
func testSkykeyListKeysDoesntShowPrivateKeys(t *testing.T, c client.Client) {
	nKeys := 3
	nExtraLines := 3
	keyStrings := make([]string, nKeys)
	keyNames := make([]string, nKeys)
	keyIDs := make([]string, nKeys)

	initSkykeyData(t, c, keyStrings, keyNames, keyIDs)

	keyListString, err := skykeyListKeys(c, false)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < nKeys; i++ {
		if !strings.Contains(keyListString, keyNames[i]) {
			t.Log(keyListString)
			t.Fatal("Missing key name!", i)
		}
		if !strings.Contains(keyListString, keyIDs[i]) {
			t.Log(keyListString)
			t.Fatal("Missing id!", i)
		}
		if strings.Contains(keyListString, keyStrings[i]) {
			t.Log(keyListString)
			t.Fatal("Found key!", i)
		}
	}
	keyList := strings.Split(keyListString, "\n")
	if len(keyList) != nKeys+nExtraLines {
		t.Fatal("Unexpected number of lines/keys", len(keyList))
	}
}

// testSkykeyListKeysAdditionalKeys tests that after creating additional keys,
// skykeyListKeys shows all key names, ids and keys
func testSkykeyListKeysAdditionalKeys(t *testing.T, c client.Client) {
	nExtraKeys := 10
	nKeys := 3 + nExtraKeys
	keyStrings := make([]string, nKeys)
	keyNames := make([]string, nKeys)
	keyIDs := make([]string, nKeys)

	initSkykeyData(t, c, keyStrings, keyNames, keyIDs)

	// Count the number of public/private keys we create.
	expectedNumPublic := nKeys - nExtraKeys // initial keys are public
	expectedNumPrivate := 0

	// Add extra keys
	for i := 0; i < nExtraKeys; i++ {
		skykeyType := pubaccesskey.TypePrivateID
		if i%2 == 0 {
			skykeyType = pubaccesskey.TypePublicID
			expectedNumPublic += 1
		} else {
			expectedNumPrivate += 1
		}

		nextName := fmt.Sprintf("extrakey-%d", i)
		keyNames = append(keyNames, nextName)
		nextSkStr, err := skykeyCreate(c, nextName, skykeyType.ToString())
		if err != nil {
			t.Fatal(err)
		}
		var nextSkykey pubaccesskey.Pubaccesskey
		err = nextSkykey.FromString(nextSkStr)
		if err != nil {
			t.Fatal(err)
		}
		keyIDs = append(keyIDs, nextSkykey.ID().ToString())
		keyStrings = append(keyStrings, nextSkStr)
	}

	// Check that all the key names and key data is there.
	keyListString, err := skykeyListKeys(c, true)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < nKeys; i++ {
		if !strings.Contains(keyListString, keyNames[i]) {
			t.Log(keyListString)
			t.Fatal("Missing key name!", i)
		}
		if !strings.Contains(keyListString, keyIDs[i]) {
			t.Log(keyListString)
			t.Fatal("Missing id!", i)
		}
		if !strings.Contains(keyListString, keyStrings[i]) {
			t.Log(keyListString)
			t.Fatal("Missing key!", i)
		}
	}

	// Check that the expected number of public/private keys were created.
	numPublic := strings.Count(keyListString, pubaccesskey.TypePublicID.ToString())
	numPrivate := strings.Count(keyListString, pubaccesskey.TypePrivateID.ToString())
	if numPublic != expectedNumPublic {
		t.Log(keyListString)
		t.Fatalf("Expected %d %s keys got %d instead", numPublic, pubaccesskey.TypePublicID.ToString(), expectedNumPublic)
	}
	if numPrivate != expectedNumPrivate {
		t.Log(keyListString)
		t.Fatalf("Expected %d %s keys got %d instead", numPrivate, pubaccesskey.TypePrivateID.ToString(), expectedNumPrivate)
	}
}

// testSkykeyListKeysAdditionalKeysDoesntShowPrivateKeys tests that after creating additional keys,
// skykeyListKeys shows all key names, ids and doesn't show private keys
func testSkykeyListKeysAdditionalKeysDoesntShowPrivateKeys(t *testing.T, c client.Client) {
	nExtraKeys := 10
	nPrevKeys := 3
	nKeys := nPrevKeys + nExtraKeys
	nExtraLines := 3
	keyStrings := make([]string, nPrevKeys)
	keyNames := make([]string, nPrevKeys)
	keyIDs := make([]string, nPrevKeys)

	initSkykeyData(t, c, keyStrings, keyNames, keyIDs)

	// Get extra keys
	for i := 0; i < nExtraKeys; i++ {
		nextName := fmt.Sprintf("extrakey-%d", i)
		keyNames = append(keyNames, nextName)
		nextSkStr, err := skykeyGet(c, nextName, "")
		if err != nil {
			t.Fatal(err)
		}
		var nextSkykey pubaccesskey.Pubaccesskey
		err = nextSkykey.FromString(nextSkStr)
		if err != nil {
			t.Fatal(err)
		}
		keyIDs = append(keyIDs, nextSkykey.ID().ToString())
		keyStrings = append(keyStrings, nextSkStr)
	}

	// Make sure key data isn't shown but otherwise the same checks pass.
	keyListString, err := skykeyListKeys(c, false)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < nKeys; i++ {
		if !strings.Contains(keyListString, keyNames[i]) {
			t.Log(keyListString)
			t.Fatal("Missing key name!", i)
		}
		if !strings.Contains(keyListString, keyIDs[i]) {
			t.Log(keyListString)
			t.Fatal("Missing id!", i)
		}
		if strings.Contains(keyListString, keyStrings[i]) {
			t.Log(keyListString)
			t.Fatal("Found key!", i)
		}
	}
	keyList := strings.Split(keyListString, "\n")
	if len(keyList) != nKeys+nExtraLines {
		t.Fatal("Unexpected number of lines/keys", len(keyList))
	}
}

// initSkykeyData initializes keyStrings, keyNames, keyIDS slices with existing Pubaccesskey data
func initSkykeyData(t *testing.T, c client.Client, keyStrings, keyNames, keyIDs []string) {
	keyName1 := "createkey1"
	keyName2 := "createkey testSkykeyGet"

	getKeyStr1, err := skykeyGet(c, keyName1, "")
	if err != nil {
		t.Fatal(err)
	}

	getKeyStr2, err := skykeyGet(c, keyName2, "")
	if err != nil {
		t.Fatal(err)
	}

	keyNames[0] = "key1"
	keyNames[1] = keyName1
	keyNames[2] = keyName2
	keyStrings[0] = testSkykeyString
	keyStrings[1] = getKeyStr1
	keyStrings[2] = getKeyStr2

	var sk pubaccesskey.Pubaccesskey
	err = sk.FromString(testSkykeyString)
	if err != nil {
		t.Fatal(err)
	}
	keyIDs[0] = sk.ID().ToString()

	err = sk.FromString(getKeyStr1)
	if err != nil {
		t.Fatal(err)
	}
	keyIDs[1] = sk.ID().ToString()

	err = sk.FromString(getKeyStr2)
	if err != nil {
		t.Fatal(err)
	}
	keyIDs[2] = sk.ID().ToString()
}
