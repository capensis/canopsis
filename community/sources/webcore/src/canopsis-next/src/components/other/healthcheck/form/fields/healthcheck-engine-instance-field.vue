<template>
  <v-layout wrap>
    <v-flex xs3>
      <v-checkbox
        v-validate
        v-field="value.enabled"
        :label="label"
        :error-messages="errors.collect(enabledFieldName)"
        :name="enabledFieldName"
        color="primary"
      />
    </v-flex>
    <v-flex xs4>
      <v-layout column>
        <v-layout align-center>
          <v-flex
            :class="{ 'text--disabled': !value.enabled }"
            xs6
          >
            {{ $t('common.minimal') }}
          </v-flex>
          <v-flex xs6>
            <c-number-field
              v-field="value.minimal"
              :error-messages="getErrorMessages(minimalFieldName)"
              :name="minimalFieldName"
              :label="$t('common.minimal')"
              :disabled="!value.enabled"
              :required="value.enabled"
              :max="+value.optimal"
              :min="0"
            />
          </v-flex>
        </v-layout>
        <v-layout align-center>
          <v-flex
            :class="{ 'text--disabled': !value.enabled }"
            xs6
          >
            {{ $t('common.optimal') }}
          </v-flex>
          <v-flex xs6>
            <c-number-field
              v-field="value.optimal"
              :error-messages="getErrorMessages(optimalFieldName)"
              :name="optimalFieldName"
              :label="$t('common.optimal')"
              :disabled="!value.enabled"
              :required="value.enabled"
              :min="+value.minimal"
              class="mt-0"
            />
          </v-flex>
        </v-layout>
      </v-layout>
    </v-flex>
    <v-flex xs9>
      <v-messages
        :value="errors.collect(name)"
        color="error"
      />
    </v-flex>
  </v-layout>
</template>

<script>
export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'engine',
    },
  },
  computed: {
    enabledFieldName() {
      return `${this.name}.enabled`;
    },

    minimalFieldName() {
      return `${this.name}.minimal`;
    },

    optimalFieldName() {
      return `${this.name}.optimal`;
    },
  },
  methods: {
    getErrorMessages(name) {
      return this.errors.collect(name, null, false)
        .map((item) => {
          const messageKey = `healthcheck.validation.${item.rule}`;

          return this.$te(messageKey) ? this.$t(messageKey) : item.msg;
        });
    },
  },
};
</script>
