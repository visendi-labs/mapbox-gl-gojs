package mapboxglgojs

import "github.com/paulmach/orb"

// TODO: make this stricter? How to separate the different options for the values?

type MapLayerPaint struct {
	// Common
	Opacity any `json:"opacity,omitempty"`

	// Fill layers
	FillColor        any    `json:"fill-color,omitempty"`
	FillOpacity      any    `json:"fill-opacity,omitempty"`
	FillOutlineColor any    `json:"fill-outline-color,omitempty"`
	FillPattern      string `json:"fill-pattern,omitempty"`

	// Line layers
	LineColor     any    `json:"line-color,omitempty"`
	LineOpacity   any    `json:"line-opacity,omitempty"`
	LineWidth     any    `json:"line-width,omitempty"`
	LineGapWidth  any    `json:"line-gap-width,omitempty"`
	LineOffset    any    `json:"line-offset,omitempty"`
	LineDashArray any    `json:"line-dasharray,omitempty"`
	LinePattern   string `json:"line-pattern,omitempty"`

	// Symbol layers
	IconOpacity   any `json:"icon-opacity,omitempty"`
	IconColor     any `json:"icon-color,omitempty"`
	IconHaloColor any `json:"icon-halo-color,omitempty"`
	IconHaloWidth any `json:"icon-halo-width,omitempty"`
	IconHaloBlur  any `json:"icon-halo-blur,omitempty"`
	TextOpacity   any `json:"text-opacity,omitempty"`
	TextColor     any `json:"text-color,omitempty"`
	TextHaloColor any `json:"text-halo-color,omitempty"`
	TextHaloWidth any `json:"text-halo-width,omitempty"`
	TextHaloBlur  any `json:"text-halo-blur,omitempty"`

	// Circle layers
	CircleRadius        any `json:"circle-radius,omitempty"`
	CircleColor         any `json:"circle-color,omitempty"`
	CircleOpacity       any `json:"circle-opacity,omitempty"`
	CircleStrokeColor   any `json:"circle-stroke-color,omitempty"`
	CircleStrokeWidth   any `json:"circle-stroke-width,omitempty"`
	CircleStrokeOpacity any `json:"circle-stroke-opacity,omitempty"`

	// Raster layers
	RasterOpacity       any `json:"raster-opacity,omitempty"`
	RasterHueRotate     any `json:"raster-hue-rotate,omitempty"`
	RasterBrightnessMin any `json:"raster-brightness-min,omitempty"`
	RasterBrightnessMax any `json:"raster-brightness-max,omitempty"`
	RasterSaturation    any `json:"raster-saturation,omitempty"`
	RasterContrast      any `json:"raster-contrast,omitempty"`
	RasterFadeDuration  any `json:"raster-fade-duration,omitempty"`

	// Hillshade layers
	HillshadeIlluminationDirection any    `json:"hillshade-illumination-direction,omitempty"`
	HillshadeIlluminationAnchor    string `json:"hillshade-illumination-anchor,omitempty"`
	HillshadeExaggeration          any    `json:"hillshade-exaggeration,omitempty"`
	HillshadeHighlightColor        any    `json:"hillshade-highlight-color,omitempty"`
	HillshadeShadowColor           any    `json:"hillshade-shadow-color,omitempty"`
	HillshadeAccentColor           any    `json:"hillshade-accent-color,omitempty"`

	// Background layers
	BackgroundColor   any    `json:"background-color,omitempty"`
	BackgroundOpacity any    `json:"background-opacity,omitempty"`
	BackgroundPattern string `json:"background-pattern,omitempty"`

	// Fill-extrusion layers (3D buildings, etc.)
	FillExtrusionColor   any    `json:"fill-extrusion-color,omitempty"`
	FillExtrusionOpacity any    `json:"fill-extrusion-opacity,omitempty"`
	FillExtrusionHeight  any    `json:"fill-extrusion-height,omitempty"`
	FillExtrusionBase    any    `json:"fill-extrusion-base,omitempty"`
	FillExtrusionPattern string `json:"fill-extrusion-pattern,omitempty"`

	// Heatmap layers
	HeatmapRadius    any `json:"heatmap-radius,omitempty"`
	HeatmapWeight    any `json:"heatmap-weight,omitempty"`
	HeatmapIntensity any `json:"heatmap-intensity,omitempty"`
	HeatmapColor     any `json:"heatmap-color,omitempty"`
	HeatmapOpacity   any `json:"heatmap-opacity,omitempty"`
}

type MapLayout struct {
	// Line properties
	LineJoin string `json:"line-join,omitempty"` // miter, bevel, round
	LineCap  string `json:"line-cap,omitempty"`  // butt, round, square

	// Symbol/Icon properties
	IconImage           string `json:"icon-image,omitempty"`
	IconSize            any    `json:"icon-size,omitempty"`
	IconRotate          any    `json:"icon-rotate,omitempty"`
	IconOffset          any    `json:"icon-offset,omitempty"`
	IconAnchor          string `json:"icon-anchor,omitempty"` // center, top, bottom, left, right, etc.
	IconAllowOverlap    bool   `json:"icon-allow-overlap,omitempty"`
	IconIgnorePlacement bool   `json:"icon-ignore-placement,omitempty"`
	IconOptional        bool   `json:"icon-optional,omitempty"`

	// Text properties
	TextField           string `json:"text-field,omitempty"`
	TextFont            any    `json:"text-font,omitempty"`
	TextSize            any    `json:"text-size,omitempty"`
	TextMaxWidth        any    `json:"text-max-width,omitempty"`
	TextLineHeight      any    `json:"text-line-height,omitempty"`
	TextLetterSpacing   any    `json:"text-letter-spacing,omitempty"`
	TextJustify         string `json:"text-justify,omitempty"`
	TextAnchor          string `json:"text-anchor,omitempty"`
	TextRotate          any    `json:"text-rotate,omitempty"`
	TextTransform       string `json:"text-transform,omitempty"`
	TextOffset          any    `json:"text-offset,omitempty"`
	TextAllowOverlap    bool   `json:"text-allow-overlap,omitempty"`
	TextIgnorePlacement bool   `json:"text-ignore-placement,omitempty"`
	TextOptional        bool   `json:"text-optional,omitempty"`

	// Visibility
	Visibility string `json:"visibility,omitempty"` // visible, none
}

type Map struct {
	// Required
	Container   string    `json:"container,omitempty"`       // HTML element id or reference
	Style       string    `json:"style,omitempty"`           // Style URL or JSON object
	Center      orb.Point `json:"center,omitempty,omitzero"` // [lng, lat]
	Zoom        float64   `json:"zoom,omitempty,omitzero"`
	AccessToken string    `json:"accessToken,omitempty"`

	// Optional view state
	Bearing          float64   `json:"bearing,omitempty,omitzero"` // Map rotation in degrees
	Pitch            float64   `json:"pitch,omitempty,omitzero"`   // Tilt in degrees
	Bounds           []float64 `json:"bounds,omitempty,omitzero"`  // [minX, minY, maxX, maxY]
	FitBoundsOptions any       `json:"fitBoundsOptions,omitempty,omitzero"`

	// Interaction controls
	Interactive     bool `json:"interactive,omitempty"`
	ScrollZoom      bool `json:"scrollZoom,omitempty"`
	DragRotate      bool `json:"dragRotate,omitempty"`
	DragPan         bool `json:"dragPan,omitempty"`
	Keyboard        bool `json:"keyboard,omitempty"`
	DoubleClickZoom bool `json:"doubleClickZoom,omitempty"`
	TouchZoomRotate bool `json:"touchZoomRotate,omitempty"`

	// Render/Performance options
	MinZoom                      float64   `json:"minZoom,omitempty,omitzero"`
	MaxZoom                      float64   `json:"maxZoom,omitempty,omitzero"`
	MaxBounds                    []float64 `json:"maxBounds,omitempty,omitzero"` // [minX, minY, maxX, maxY]
	PreserveDrawingBuffer        bool      `json:"preserveDrawingBuffer,omitempty"`
	Antialias                    bool      `json:"antialias,omitempty"`
	TrackResize                  bool      `json:"trackResize,omitempty"`
	FailIfMajorPerformanceCaveat bool      `json:"failIfMajorPerformanceCaveat,omitempty"`

	// Locale and attribution
	Locale             map[string]string `json:"locale,omitempty"`
	AttributionControl bool              `json:"attributionControl,omitempty"`

	// Misc
	Hash bool `json:"hash,omitempty"` // Track position in URL hash

	// Your custom config
	Config MapConfig `json:"config,omitempty,omitzero"`
}

type MapConfig struct {
	Basemap BasemapConfig `json:"basemap,omitempty,omitzero"`
}

type BasemapConfig struct {
	Theme                  string `json:"theme,omitempty"`
	Show3DObjects          bool   `json:"show3dObjects,omitempty"`
	ShowPointOfInterest    bool   `json:"showPointOfInterestLabels,omitempty"`
	ShowPlaceLabels        bool   `json:"showPlaceLabels,omitempty"`
	ShowRoadLabels         bool   `json:"showRoadLabels,omitempty"`
	ShowTransitLabels      bool   `json:"showTransitLabels,omitempty"`
	ShowAdminBoundaries    bool   `json:"showAdminBoundaries,omitempty"`
	ShowBuildingExtrusions bool   `json:"showBuildingExtrusions,omitempty"`
	ShowRoads              bool   `json:"showRoads,omitempty"`
	ShowTransit            bool   `json:"showTransit,omitempty"`
	ShowTerrain            bool   `json:"showTerrain,omitempty"`
	ShowWater              bool   `json:"showWater,omitempty"`
}
