/// <reference path="../pb_data/types.d.ts" />
migrate((app) => {
  const collection = app.findCollectionByNameOrId("pbc_3598190544")

  // update collection data
  unmarshal({
    "name": "galleries"
  }, collection)

  return app.save(collection)
}, (app) => {
  const collection = app.findCollectionByNameOrId("pbc_3598190544")

  // update collection data
  unmarshal({
    "name": "gallery"
  }, collection)

  return app.save(collection)
})
