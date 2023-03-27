<template lang="pug">
  v-layout.my-2(v-if="!form", justify-center)
    v-progress-circular(color="primary", indeterminate)
  v-flex(v-else)
    v-layout(row)
      v-flex.text-xs-center
        div.title {{ $t('userInterface.title') }}
    v-form(@submit.prevent="submit")
      user-interface-form(v-model="form", ref="userForm")
      template(v-if="!disabled")
        v-divider.mt-3
        v-layout.mt-3(row, justify-end)
          v-btn(flat, @click="reset") {{ $t('common.cancel') }}
          v-btn.primary(
            :disabled="isDisabled",
            :loading="submitting",
            type="submit"
          ) {{ $t('common.submit') }}
</template>

<script>
import { VALIDATION_DELAY } from '@/constants';

import { userInterfaceToForm } from '@/helpers/forms/user-interface';
import { getFileDataUrlContent } from '@/helpers/file/file-select';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { submittableMixinCreator } from '@/mixins/submittable';

import UserInterfaceForm from './form/user-interface-form.vue';

export default {
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    UserInterfaceForm,
  },
  mixins: [
    entitiesInfoMixin,
    submittableMixinCreator(),
  ],
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      form: null,
    };
  },
  async mounted() {
    await this.fetchAppInfo();

    this.form = userInterfaceToForm(this.appInfo);
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        const data = { ...this.form };

        if (this.form.logo) {
          data.logo = await getFileDataUrlContent(this.form.logo);
        }

        await this.updateUserInterface({ data });
        await this.fetchAppInfo();

        this.setTitle();

        this.$popups.success({ text: this.$t('success.default') });

        this.reset();
      }
    },

    reset() {
      this.$refs.userForm.reset();

      this.form = userInterfaceToForm(this.appInfo);
    },
  },
};
</script>
