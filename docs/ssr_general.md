# Serverside Rendering + Mapbox + Golang

Mapbox-GL-GOJS can be used in a SSR web app context. The following examples will dive deeper into this.

Mapbox + SSR don't go that well together. It can force you to have to handle lots of map/layer/source logic client-side in JS/TS. This can cause all sorts of headaches.

Keeping your map source of truth serverside comes with a ton of benefits. Here are some:
- With all of your app logic serverside, it's no more than right to keep your map logic there too
    - DRY, no duplicate code to handle the same logic both clientside and serverside
    - One language for all logic
    - Everyday server things like saving to session easily integrated with map operations
    - Server can cache or do other smart things when fetching data to show
- No tricky SSR templating, injecting some serverside data into some function call in a JS/HTML template. 

Some cons:
- Might lead to more data/traffic being sent from server
    - Look into gzipping responses, and adding cache headers so that the user can cache responses. Have a look at the examples. 