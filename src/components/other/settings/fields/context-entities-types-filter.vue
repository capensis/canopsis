<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.defaultSortColumn')}}
    v-container
      v-select(:value="value", @input="input", :items="entitiesTypes")
</template>

<script>
export default {
  props: {
    value: String,
  },
  data() {
    return {
      entitiesTypes: ['', 'component', 'connector', 'resource', 'watcher'],
    };
  },
  methods: {
    input(value) {
      this.$emit('input', value);
    },
  },
};
</script>
