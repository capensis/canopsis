<template lang="pug">
  div
    slot(v-if="isPatternsEmpty", name="no-data")
      v-alert(
        :value="true",
        type="info"
      ) {{ disabled ? $t('patternsList.noDataDisabled') : $t('patternsList.noData') }}
    v-layout(
      v-for="(pattern, index) in patterns",
      :key="`${$options.filters.json(pattern)}${index}`",
      row,
      wrap,
      align-center
    )
      v-flex(:class="disabled ? 'xs12' : 'xs11'")
        v-layout
          pattern-information.ma-3.ml-0 {{ $t('common.and') }}
          v-flex(xs12)
            v-textarea(
              :value="pattern | json",
              rows="7",
              no-resize,
              readonly,
              disabled
            )
        v-layout(v-if="index !== patterns.length - 1", justify-center)
          span.text-uppercase.operator-chip {{ $t('common.or') }}
      v-flex.text-xs-center(v-if="!disabled", xs1)
        div
          v-btn(icon, @click="showEditPatternModal(index)")
            v-icon edit
        div
          v-btn(color="error", icon, @click="showRemovePatternModal(index)")
            v-icon delete
    v-btn.mx-0(v-if="!disabled", color="primary", @click="showCreatePatternModal") {{ $t('common.add') }}
    v-alert(
      :value="errors.has(name)",
      type="error",
      transition="fade-transition"
    )
      span(v-for="error in errors.collect(name)", :key="error") {{ error }}
    v-alert(
      v-model="countAlertShown",
      type="warning",
      transition="fade-transition",
      dismissible
    )
      span {{ countAlertMessage }}
</template>

<script>
import { Validator } from 'vee-validate';

import { MODALS, EVENT_FILTER_OPERATORS } from '@/constants';

import { formArrayMixin } from '@/mixins/form';

import PatternInformation from '@/components/other/pattern/pattern-information.vue';

export default {
  $_veeValidate: {
    value() {
      return this.patterns;
    },

    name() {
      return this.name;
    },
  },
  inject: {
    $validator: {
      default() {
        return new Validator();
      },
    },
    $checkEntitiesCountForPatternsByType: {
      default() {
        return null;
      },
    },
  },
  components: {
    PatternInformation,
  },
  mixins: [formArrayMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Array,
      default: () => [],
    },
    operators: {
      type: Array,
      default: () => EVENT_FILTER_OPERATORS,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'patterns',
    },
    type: {
      type: String,
      required: false,
    },
  },
  data() {
    return {
      countAlertShown: false,
      countAlertMessage: '',
    };
  },
  computed: {
    isPatternsEmpty() {
      return !this.patterns || !this.patterns.length;
    },
  },
  methods: {
    showCreatePatternModal() {
      this.$modals.show({
        name: MODALS.createPattern,
        config: {
          operators: this.operators,
          action: async (pattern) => {
            const patterns = this.addItemIntoArray(pattern);

            this.errors.remove(this.name);

            await this.checkEntitiesCount(patterns);
          },
        },
      });
    },

    showEditPatternModal(index) {
      this.$modals.show({
        name: MODALS.createPattern,
        config: {
          pattern: this.patterns[index],
          operators: this.operators,
          action: async (pattern) => {
            const patterns = this.updateItemInArray(index, pattern);

            this.errors.remove(this.name);

            await this.checkEntitiesCount(patterns);
          },
        },
      });
    },

    showRemovePatternModal(index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: async () => {
            const patterns = this.removeItemFromArray(index);

            this.errors.remove(this.name);

            await this.checkEntitiesCount(patterns);
          },
        },
      });
    },

    async checkEntitiesCount(patterns = []) {
      if (!this.$checkEntitiesCountForPatternsByType || !this.type || !patterns.length) {
        this.countAlertShown = false;

        return;
      }

      try {
        const {
          over_limit: overLimit,
          total_count: totalCount,
        } = await this.$checkEntitiesCountForPatternsByType(this.type, patterns);

        if (overLimit) {
          this.countAlertMessage = this.$t('entitiesCountAlerts.patterns.countOverLimit', { count: totalCount });
          this.countAlertShown = true;

          return;
        }

        this.countAlertShown = false;
      } catch (err) {
        this.countAlertMessage = this.$t('entitiesCountAlerts.patterns.countRequestError');
        this.countAlertShown = true;
      }
    },
  },
};
</script>
