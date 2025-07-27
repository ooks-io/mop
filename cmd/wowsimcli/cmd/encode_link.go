package cmd

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/wowsims/mop/sim/core/proto"
	"google.golang.org/protobuf/encoding/protojson"
	goproto "google.golang.org/protobuf/proto"
)

var encodeLinkCmd = &cobra.Command{
	Use:   "encodelink [json-file]",
	Short: "encode simulation settings to wowsims link/url",
	Long:  "encode simulation settings from JSON file to wowsims shareable link/url",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return encodeLink(args[0])
	},
}

func encodeLink(jsonFile string) error {
	// read json input file
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		return fmt.Errorf("failed to read JSON file %q: %w", jsonFile, err)
	}

	// try to determine sim type by checking if it's a raid or individual sim
	var settings goproto.Message
	var urlPath string
	baseUrl := "https://www.wowsims.com/mop"

	// first try individual sim
	individualSettings := &proto.IndividualSimSettings{}
	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, individualSettings)
	if err == nil && individualSettings.Player != nil {
		// it's an individual sim
		settings = individualSettings
		classStr, specStr := getClassAndSpecFromPlayer(individualSettings.Player)
		if classStr != "" && specStr != "" {
			urlPath = fmt.Sprintf("/%s/%s/", classStr, specStr)
		} else {
			urlPath = "/individual/" // fallback
		}
	} else {
		// try raid sim
		raidSettings := &proto.RaidSimSettings{}
		err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, raidSettings)
		if err != nil {
			return fmt.Errorf("failed to unmarshal JSON as either individual or raid sim: %w", err)
		}
		settings = raidSettings
		urlPath = "/raid/"
	}

	// marshal protobuf to binary
	protoBytes, err := goproto.Marshal(settings)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf to binary: %w", err)
	}

	// compress with zlib
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write(protoBytes); err != nil {
		return fmt.Errorf("failed to write to zlib compressor: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("failed to close zlib compressor: %w", err)
	}

	// base64 encode the compressed data
	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	// generate the shareable URL
	shareableUrl := fmt.Sprintf("%s%s#%s", baseUrl, urlPath, encoded)
	fmt.Println(shareableUrl)

	return nil
}

// getClassAndSpecFromPlayer extracts class and spec strings from Player protobuf
func getClassAndSpecFromPlayer(player *proto.Player) (string, string) {
	specMap := map[string][2]string{
		"*proto.Player_BloodDeathKnight":  {"death_knight", "blood"},
		"*proto.Player_FrostDeathKnight":  {"death_knight", "frost"},
		"*proto.Player_UnholyDeathKnight": {"death_knight", "unholy"},

		"*proto.Player_BalanceDruid":     {"druid", "balance"},
		"*proto.Player_FeralDruid":       {"druid", "feral"},
		"*proto.Player_GuardianDruid":    {"druid", "guardian"},
		"*proto.Player_RestorationDruid": {"druid", "restoration"},

		"*proto.Player_BeastMasteryHunter": {"hunter", "beast_mastery"},
		"*proto.Player_MarksmanshipHunter": {"hunter", "marksmanship"},
		"*proto.Player_SurvivalHunter":     {"hunter", "survival"},

		"*proto.Player_ArcaneMage": {"mage", "arcane"},
		"*proto.Player_FireMage":   {"mage", "fire"},
		"*proto.Player_FrostMage":  {"mage", "frost"},

		"*proto.Player_BrewmasterMonk": {"monk", "brewmaster"},
		"*proto.Player_MistweaverMonk": {"monk", "mistweaver"},
		"*proto.Player_WindwalkerMonk": {"monk", "windwalker"},

		"*proto.Player_HolyPaladin":        {"paladin", "holy"},
		"*proto.Player_ProtectionPaladin":  {"paladin", "protection"},
		"*proto.Player_RetributionPaladin": {"paladin", "retribution"},

		"*proto.Player_DisciplinePriest": {"priest", "discipline"},
		"*proto.Player_HolyPriest":       {"priest", "holy"},
		"*proto.Player_ShadowPriest":     {"priest", "shadow"},

		"*proto.Player_AssassinationRogue": {"rogue", "assassination"},
		"*proto.Player_CombatRogue":        {"rogue", "combat"},
		"*proto.Player_SubtletyRogue":      {"rogue", "subtlety"},

		"*proto.Player_ElementalShaman":   {"shaman", "elemental"},
		"*proto.Player_EnhancementShaman": {"shaman", "enhancement"},
		"*proto.Player_RestorationShaman": {"shaman", "restoration"},

		"*proto.Player_AfflictionWarlock":  {"warlock", "affliction"},
		"*proto.Player_DemonologyWarlock":  {"warlock", "demonology"},
		"*proto.Player_DestructionWarlock": {"warlock", "destruction"},

		"*proto.Player_ArmsWarrior":       {"warrior", "arms"},
		"*proto.Player_FuryWarrior":       {"warrior", "fury"},
		"*proto.Player_ProtectionWarrior": {"warrior", "protection"},
	}

	typeName := fmt.Sprintf("%T", player.Spec)
	if mapping, exists := specMap[typeName]; exists {
		return mapping[0], mapping[1]
	}
	return "", ""
}
