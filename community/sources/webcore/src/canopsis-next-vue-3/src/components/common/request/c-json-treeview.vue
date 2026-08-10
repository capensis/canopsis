<template>
  <pre class="json-treeview">{{ formattedJson }}
</pre>
</template>

<script>
export default {
  props: {
    json: {
      type: String,
      default: '',
    },
    jsonObject: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    parsedJson() {
      return this.json ? JSON.parse(this.json) : this.jsonObject;
    },

    formattedJson() {
      return JSON.stringify(this.parsedJson, null, 2);
    },
  },
};
</script>

<style lang="scss">
.json-treeview {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
