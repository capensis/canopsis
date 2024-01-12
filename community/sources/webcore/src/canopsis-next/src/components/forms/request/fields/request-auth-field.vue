<template>
  <v-layout>
    <v-flex xs4>
      <c-enabled-field
        v-field="value.enabled"
        :label="$t('scenario.withAuth')"
        :disabled="disabled"
      />
    </v-flex>
    <template v-if="value.enabled">
      <v-flex
        class="mr-3"
        xs4
      >
        <c-name-field
          v-field="value.username"
          :label="$t('common.username')"
          :name="`${name}.username`"
          :disabled="disabled"
        />
      </v-flex>
      <v-flex xs4>
        <c-password-field
          v-field="value.password"
          :name="`${name}.password`"
          :disabled="disabled"
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
  },
};
</script>
