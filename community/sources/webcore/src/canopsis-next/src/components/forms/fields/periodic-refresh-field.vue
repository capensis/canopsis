<template>
  <v-layout
    class="mb-3"
    align-top
  >
    <v-flex xs5>
      <c-enabled-field
        v-field="periodicRefresh.enabled"
        :label="label"
        hide-details
        @input="validateDuration"
      />
    </v-flex>
    <v-flex xs7>
      <c-duration-field
        v-field="periodicRefresh"
        :disabled="!periodicRefresh.enabled"
        :required="periodicRefresh.enabled"
        :name="name"
        :min="1"
      />
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'periodicRefresh',
    event: 'input',
  },
  props: {
    periodicRefresh: {
      type: Object,
      default: () => ({}),
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'periodic_refresh',
    },
  },
  methods: {
    validateDuration() {
      this.$nextTick(() => this.$validator.validate(`${this.name}.value`));
    },
  },
};
</script>
