package main

import (
	"fmt"
	"log/slog"

	"github.com/redhatinsights/rhc/internal/conf"
	"github.com/redhatinsights/rhc/internal/features"
	"github.com/urfave/cli/v2"
)

// beforeEnableFeaturesAction is called before enableFeaturesAction
func beforeEnableFeaturesAction(ctx *cli.Context) error {
	slog.Debug("Command 'rhc configure features enable' started")

	err := setupFormatOption(ctx)
	if err != nil {
		return err
	}

	configureUI(ctx)

	if ctx.Args().Len() == 0 {
		err = fmt.Errorf("error: required argument 'FEATURE' is missing")
		return cli.Exit(err, ExitCodeDataErr)
	}

	return nil
}

// enableFeaturesAction tries to enables features
func enableFeaturesAction(ctx *cli.Context) error {
	args := ctx.Args()
	enabledFeatures := args.Slice()
	var disabledFeatures []string

	// First, validate provided features
	err := features.ValidateSelectedFeatures(&enabledFeatures, &disabledFeatures)
	if err != nil {
		slog.Warn(fmt.Sprintf("Failed to validate features: %s", err.Error()))
		return cli.Exit(err, ExitCodeDataErr)
	}

	// Log the enabled features
	for _, feature := range enabledFeatures {
		slog.Debug(fmt.Sprintf("Trying to enable feature: %s", feature))
	}

	// Then check that enabled features are not in conflict with the configuration file
	enabledFeatures, disabledFeatures, err = features.ConsolidateSelectedFeatures(
		&conf.Config,
		enabledFeatures,
		disabledFeatures,
	)
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to consolidate selected features: %s", err.Error()))
		return cli.Exit(err, ExitCodeDataErr)
	}

	mapFeatures := features.MapKnownFeatureIds()
	slog.Debug("Features enabled using CLI arguments consolidated with configuration file")

	// Try to enable the features
	for _, featureId := range enabledFeatures {
		feature := mapFeatures[featureId]
		if feature == nil {
			slog.Warn(fmt.Sprintf("feature '%s' not in the map of known features", featureId))
			continue
		}
		if feature.IsEnabledFunc() {
			slog.Debug(fmt.Sprintf("Feature '%s' is already enabled", featureId))
			continue
		}
		slog.Debug(fmt.Sprintf("Feature '%s' will be enabled", featureId))
		err = feature.EnableFunc(ctx)
		if err != nil {
			slog.Warn(fmt.Sprintf("failed to enable feature '%s': %s", featureId, err.Error()))
		}
	}

	// Should we do anything with disabled features?
	for _, featureId := range disabledFeatures {
		feature := mapFeatures[featureId]
		if feature == nil {
			slog.Warn(fmt.Sprintf("feature '%s' not in the map of known features", featureId))
			continue
		}
		if feature.IsEnabledFunc() == false {
			slog.Debug(fmt.Sprintf("Feature '%s' is already disabled", featureId))
			continue
		}
		//slog.Debug(fmt.Sprintf("Feature '%s' will be disabled", featureId))
		//err = feature.DisableFunc(ctx)
		//if err != nil {
		//	slog.Warn(fmt.Sprintf("failed to disable feature '%s': %s", featureId, err.Error()))
		//}
	}

	return nil
}

func beforeDisableFeaturesAction(ctx *cli.Context) error {
	slog.Debug("Command 'rhc configure features disable' started")

	err := setupFormatOption(ctx)
	if err != nil {
		return err
	}

	configureUI(ctx)

	if ctx.Args().Len() == 0 {
		err = fmt.Errorf("error: required argument 'FEATURE' is missing")
		return cli.Exit(err, ExitCodeDataErr)
	}

	return nil
}

// disableFeaturesAction disables features
func disableFeaturesAction(ctx *cli.Context) error {
	args := ctx.Args()
	for _, arg := range args.Slice() {
		slog.Info(fmt.Sprintf("Disabling feature: %s", arg))
	}
	return nil
}

// beforeShowFeaturesAction is called before showFeaturesAction
func beforeShowFeaturesAction(ctx *cli.Context) error {
	slog.Debug("Command 'rhc configure features show' started")

	err := setupFormatOption(ctx)
	if err != nil {
		return err
	}

	configureUI(ctx)

	return checkForUnknownArgs(ctx)
}

// showFeaturesAction shows features
func showFeaturesAction(ctx *cli.Context) error {
	// content
	if features.ContentFeature.IsEnabledFunc() {
		slog.Debug("Content is enabled")
	} else {
		slog.Debug("Content is disabled")
	}
	contentEnabled := features.ContentFeature.IsEnabledInConf()
	if contentEnabled != nil {
		if *contentEnabled {
			slog.Debug("Content is enabled in configuration file")
		} else {
			slog.Debug("Content is disabled in configuration file")
		}
	} else {
		slog.Debug("Content is not configured in configuration file")
	}

	// analytics
	if features.AnalyticsFeature.IsEnabledFunc() {
		slog.Debug("Analytics is enabled")
	} else {
		slog.Debug("Analytics is disabled")
	}
	analyticsEnabled := features.AnalyticsFeature.IsEnabledInConf()
	if analyticsEnabled != nil {
		if *analyticsEnabled {
			slog.Debug("Analytics is enabled in configuration file")
		} else {
			slog.Debug("Analytics is disabled in configuration file")
		}
	} else {
		slog.Debug("Analytics is not configured in configuration file")
	}

	// management
	if features.ManagementFeature.IsEnabledFunc() {
		slog.Debug("Management is enabled")
	} else {
		slog.Debug("Management is disabled")
	}
	managementEnabled := features.ManagementFeature.IsEnabledInConf()
	if managementEnabled != nil {
		if *managementEnabled {
			slog.Debug("Management is enabled in configuration file")
		} else {
			slog.Debug("Management is disabled in configuration file")
		}
	} else {
		slog.Debug("Management is not configured in configuration file")
	}

	return nil
}
