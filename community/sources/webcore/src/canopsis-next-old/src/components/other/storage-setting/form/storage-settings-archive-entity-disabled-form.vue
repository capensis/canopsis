<template>
  <v-checkbox
    v-field="form.with_dependencies"
    color="primary"
  >
    <template #label="">
      {{ $t('storageSetting.entity.archiveDependencies') }}
      <c-help-icon
        :text="$t('storageSetting.entity.archiveDependenciesHelp')"
        icon-class="ml-2"
        color="info"
        max-width="300"
        top="top"
      />
    </template>
  </v-checkbox>
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>
