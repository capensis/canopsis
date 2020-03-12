<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span Edit dynamic info template
      template(slot="text")
        div
          v-text-field(
            v-model="form.title",
            v-validate="'required'",
            :error-messages="errors.collect('title')",
            label="Title",
            name="title"
          )
          h3 Values
          v-layout(v-for="(value, index) in form.values", :key="value.key", row, justify-space-between)
            v-flex(xs11)
              v-text-field(
                v-model="value.name",
                placeholder="Value"
              )
            v-flex(xs1)
              v-btn(
                color="error",
                icon,
                @click="deleteValue(index)"
              )
                v-icon delete
          v-btn.primary(@click="showAddValueModal") Add new value
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import uid from '@/helpers/uid';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create dynamic info's information
 */
export default {
  name: MODALS.createDynamicInfoTemplate,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    return {
      form: {
        id: uid(),
        title: 'Test template',
        values: [
          { key: uid(), name: 'attr_dynamique1' },
          { key: uid(), name: 'attr_dynamique2' },
          { key: uid(), name: 'attr_statique' },
        ],
      },
    };
  },
  methods: {
    showAddValueModal() {
      this.form.values.push({ id: uid(), name: '' });
    },
    deleteValue(index) {
      this.form.values = this.form.values.filter((v, i) => i !== index);
    },
    submit() {

    },
  },
};
</script>
