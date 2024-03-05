<template>
  <v-radio-group
    v-field="value"
    hide-details
    row
  >
    <v-radio
      :label="$t('view.duplicateAsPrivate')"
      :value="true"
      color="primary"
    />
    <v-radio
      :label="$t('view.duplicateAsRegular')"
      :value="false"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Boolean,
      required: true,
    },
  },
};
</script>
