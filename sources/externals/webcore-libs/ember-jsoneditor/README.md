[Ember-JSONEditor](https://github.com/Glavin001/ember-jsoneditor)
==========

> Ember component for [JSONEditor](https://github.com/josdejong/jsoneditor/) to view, edit and format JSON.

**Live demo:** http://glavin001.github.io/ember-jsoneditor/dist/

---

## Usage

### Installation

Follow the [installation instructions for JSONEditor](https://github.com/josdejong/jsoneditor/#install)
then install [ember-jsoneditor](https://github.com/Glavin001/ember-jsoneditor).

```bash
bower install --save ember-jsoneditor
```

### Ember Component

```handlebars
{{json-editor json=model mode=controller.mode name=controller.name}}
```

## Developing

After cloning repository, install library dependencies.

```bash
npm install
bower install
```

Then build with `grunt`.

```bash
grunt serve
```
