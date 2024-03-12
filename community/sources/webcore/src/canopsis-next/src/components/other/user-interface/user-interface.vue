<template>
  <v-layout
    v-if="!form"
    class="my-2"
    justify-center
  >
    <v-progress-circular
      color="primary"
      indeterminate
    />
  </v-layout>
  <v-flex v-else>
    <v-layout>
      <v-flex class="text-center">
        <div class="text-h6">
          {{ $t('userInterface.title') }}
        </div>
      </v-flex>
    </v-layout>
    <v-form @submit.prevent="submit">
      <user-interface-form
        v-model="form"
        ref="userForm"
      />
      <template v-if="!disabled">
        <v-divider class="mt-3" />
        <v-layout
          class="mt-3"
          justify-end
        >
          <v-btn
            text
            @click="reset"
          >
            {{ $t('common.cancel') }}
          </v-btn>
          <v-btn
            :disabled="isDisabled"
            :loading="submitting"
            class="primary"
            type="submit"
          >
            {{ $t('common.submit') }}
          </v-btn>
        </v-layout>
      </template>
    </v-form>
  </v-flex>
</template>

<script>
import { VALIDATION_DELAY } from '@/constants';

import { userInterfaceToForm } from '@/helpers/entities/user-interface/form';
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
      form: userInterfaceToForm(),
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
