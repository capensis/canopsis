<template>
  <v-layout class="gap-2">
    <v-flex v-if="!required" xs4>
      <c-enabled-field
        v-field="value.enabled"
        :label="$t('scenario.withAuth')"
        :disabled="disabled"
      />
    </v-flex>
    <template v-if="value.enabled || required">
      <v-flex :xs4="!required">
        <c-name-field
          v-field="value.username"
          :label="$t('common.username')"
          :name="`${name}.username`"
          :disabled="disabled"
          :required="required"
        />
      </v-flex>
      <v-flex :xs4="!required">
        <c-password-field
          v-field="value.password"
          :name="`${name}.password`"
          :disabled="disabled"
          :required="required"
        />
      </v-flex>
    </template>
  </v-layout>
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'auth',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
