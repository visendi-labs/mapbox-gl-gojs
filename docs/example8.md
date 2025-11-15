# Serverside Popups
[Full code](https://github.com/visendi-labs/mapbox-gl-gojs/tree/main/docs/example8) (run with `go run main.go` from the `golang` folder)

[](example8/wasm/index.html ':include :type=iframe width=100% height=500px')


`index.html` Parsed by Go's [html/template](https://pkg.go.dev/html/template) or another templating tool

[filename](/example8/golang/index.html ':include :type=code :fragment=demo')



`map.go` 

[filename](/example8/common/map.go ':include :type=code :fragment=demo')

`main.go`
[filename](/example8/golang/main.go ':include :type=code :fragment=demo')

# Comments 


The resulting JS from the call to `mbgojs.NewMapOnEventLayer("click", "points", mbgojs.NewHtmxAjax(...` will result in the following (key here is the `mbgojs.HtmxAjaxContextEventValuesFull`, making sure attributes of `e` are added to the `htmx.ajax` values):
```javascript
map.on("click", "points", (e) => {
    htmx.ajax("GET", "/popup", {
        "values": {
            "eventType": e.type,
            "featureId": e.features[0].id,
            "lat": e.lngLat.lat,
            "layerId": e.features[0].layer.id,
            "layerType": e.features[0].layer.type,
            "lng": e.lngLat.lng,
            "sourceId": e.features[0].source,
            "type": e.features[0].type,
            "x": e.point.x,
            "y": e.point.y,
        },
        "swap": "beforeend",
        "target": "body",
    });
});
```
