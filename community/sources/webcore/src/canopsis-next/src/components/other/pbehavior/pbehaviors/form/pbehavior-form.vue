<template>
  <pbehavior-general-form
    v-if="noPattern"
    v-field="form"
    :no-enabled="noEnabled"
    :no-comments="noComments"
    :with-start-on-trigger="withStartOnTrigger"
    :name-label="nameLabel"
    :name-tooltip="nameTooltip"
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
      <pbehavior-general-form
        v-field="form"
        ref="general"
        :no-enabled="noEnabled"
        :no-comments="noComments"
        :with-start-on-trigger="withStartOnTrigger"
        :name-label="nameLabel"
        :name-tooltip="nameTooltip"
      />
    </v-tab-item>
    <v-tab-item eager>
      <pbehavior-patterns-form
        v-field="form.patterns"
        ref="patterns"
      />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import PbehaviorGeneralForm from './pbehavior-general-form.vue';
import PbehaviorPatternsForm from './pbehavior-patterns-form.vue';

export default {
  inject: ['$validator'],
  components: {
    PbehaviorGeneralForm,
    PbehaviorPatternsForm,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    noPattern: {
      type: Boolean,
      default: false,
    },
    noEnabled: {
      type: Boolean,
      default: false,
    },
    noComments: {
      type: Boolean,
      default: false,
    },
    withStartOnTrigger: {
      type: Boolean,
      default: false,
    },
    nameLabel: {
      type: String,
      required: false,
    },
    nameTooltip: {
      type: String,
      required: false,
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
      this.unwatchGeneralTabErrors = this.$watch(() => this.$refs.general.hasAnyError, (value) => {
        this.hasGeneralError = value;
      });

      this.unwatchPatternsTabErrors = this.$watch(() => this.$refs.patterns.hasAnyError, (value) => {
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
