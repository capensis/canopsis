<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createDynamicInfo.create.title') }}
    v-card-text
      v-form
        dynamic-info-form(v-model="form")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
      v-btn(color="primary", @click="submit") {{ $t('common.submit') }}
</template>

<script>
import { isEmpty } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import DynamicInfoForm from '@/components/other/dynamic-info/form/dynamic-info-form.vue';

/**
 * Modal to create dynamic information
 */
export default {
  name: MODALS.createDynamicInfo,
  $_veeValidate: {
    validator: 'new',
  },
  components: { DynamicInfoForm },
  mixins: [modalInnerMixin],
  data() {
    return {
      stepper: 0,
      form: {
        general: {
          _id: '',
          name: '',
          description: '',
        },
        infos: [],
        patterns: {
          alarm_patterns: [],
          entity_patterns: [],
        },
      },
    };
  },
  created() {
    this.$validator.attach({
      name: 'pattern',
      rules: 'required:true',
      getter: () => !isEmpty(this.form.pattern),
      context: () => this,
    });
  },
  methods: {
    async submit() {
      try {
        const isValid = await this.$validator.validateAll();

        if (isValid) {
          // TODO: Prepare data object
          const data = {};

          if (this.config.action) {
            await this.config.action(data);
          }

          this.$modals.hide();
        }
      } catch (err) {
        this.$popups.error({ text: err.description });
      }
    },
  },
};
</script>
