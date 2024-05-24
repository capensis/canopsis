<template>
  <div>
    <remediation-instruction-general-form
      v-if="noPattern"
      v-field="form"
      :disabled="disabled"
      :is-new="isNew"
      :required-approve="requiredApprove"
    />
    <v-tabs
      v-else
      slider-color="primary"
      centered
    >
      <v-tab :class="{ 'error--text': hasGeneralError }">
        {{ $t('common.general') }}
      </v-tab>
      <v-tab :class="{ 'error--text': hasPatternsError }">
        {{ $tc('common.pattern', 2) }}
      </v-tab>

      <v-tab-item eager>
        <remediation-instruction-general-form
          v-field="form"
          ref="general"
          :disabled="disabled"
          :is-new="isNew"
          :required-approve="requiredApprove"
          class="mt-3"
        />
      </v-tab-item>
      <v-tab-item eager>
        <remediation-instruction-patterns-form
          v-field="form.patterns"
          ref="patterns"
          class="mt-3"
        />
      </v-tab-item>
    </v-tabs>
  </div>
</template>

<script>
import RemediationInstructionGeneralForm from './remediation-instruction-general-form.vue';
import RemediationInstructionPatternsForm from './remediation-instruction-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    RemediationInstructionGeneralForm,
    RemediationInstructionPatternsForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    noPattern: {
      type: Boolean,
      default: false,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    disabledCommon: {
      type: Boolean,
      default: false,
    },
    isNew: {
      type: Boolean,
      default: false,
    },
    requiredApprove: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      hasGeneralError: false,
      hasPatternsError: false,
    };
  },
  watch: {
    noPattern: {
      handler(noPattern) {
        if (noPattern) {
          this.unwatchTabsErrors();
        } else {
          this.$nextTick(this.watchTabsErrors);
        }
      },
      immediate: true,
    },
  },
  methods: {
    watchTabsErrors() {
      this.unwatchGeneralTabErrors = this.$watch(() => this.$refs.general?.hasAnyError, (value) => {
        this.hasGeneralError = value;
      });

      this.unwatchPatternsTabErrors = this.$watch(() => this.$refs.patterns?.hasAnyError, (value) => {
        this.hasPatternsError = value;
      });
    },

    unwatchTabsErrors() {
      this.unwatchGeneralTabErrors?.();
      this.unwatchPatternsTabErrors?.();
    },
  },
};
</script>
