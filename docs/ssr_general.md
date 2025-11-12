# Serverside Rendering + Mapbox + Golang

Mapbox-GL-GOJS can be used in a SSR web app context. The following sections will dive deeper into this.

Mapbox + SSR don't go very well together. Using Mapbox can force you to handle lots of map logic clientside in JS/TS. This can cause all sorts of headaches.

Keeping your map source of truth serverside comes with a ton of benefits. Here are some:
- With all of your app logic serverside, it's ideal to keep your map logic there too
    - DRY, no duplicate code to handle the same logic both clientside and serverside
    - One language for all logic
    - Server things like saving to session runs smoothly
    - Server can cache or do other smart things when fetching data to show
- No tricky SSR templating, injecting some serverside data into some function call in a JS/HTML template. 

Some cons:
- Might lead to more data/traffic being sent from server
    - Look into gzipping responses, and adding cache headers so that the user can cache responses