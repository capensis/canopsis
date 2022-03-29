<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ title }}
      template(#text="")
        c-patterns-field(v-model="form", with-alarm, with-entity)
        c-collapse-panel(color="grey")
          template(#header="")
            span.white--text {{ $t('remediationPatterns.tabs.pbehaviorTypes.title') }}
          v-card
            v-card-text
              remediation-patterns-pbehavior-types-form(v-model="form")
      template(#actions="")
        v-btn(
          :disabled="submitting",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, PATTERNS_FIELDS } from '@/constants';

import { filterPatternsToForm, formFilterToPatterns } from '@/helpers/forms/filter';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import RemediationPatternsPbehaviorTypesForm
  from '@/components/other/remediation/patterns/remediation-patterns-pbehavior-types-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.remediationPatterns,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    RemediationPatternsPbehaviorTypesForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { instruction } = this.modal.config;

    return {
      form: {
        ...filterPatternsToForm(instruction, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
        active_on_pbh: instruction.active_on_pbh ?? [],
        disabled_on_pbh: instruction.disabled_on_pbh ?? [],
      },
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.patterns.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action({
            ...formFilterToPatterns(this.form, [PATTERNS_FIELDS.alarm, PATTERNS_FIELDS.entity]),
            active_on_pbh: this.form.active_on_pbh,
            disabled_on_pbh: this.form.disabled_on_pbh,
          });
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
