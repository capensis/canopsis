<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(#title="")
        span {{ $t('modals.infoPopupSetting.title') }}
      template(#text="")
        v-layout(justify-end)
          v-btn(
            icon,
            fab,
            small,
            color="secondary",
            @click="addPopup"
          )
            v-icon add
        v-layout(column)
          v-card.my-1(
            v-for="(popup, index) in form.popups",
            :key="index",
            color="secondary white--text",
            flat
          )
            v-card-title
              v-layout(justify-space-between)
                div {{ $t('modals.infoPopupSetting.column') }}: {{ popup.column }}
                div
                  v-btn(
                    icon,
                    small,
                    @click="deletePopup(index)"
                  )
                    v-icon(color="error") delete
                  v-btn(
                    icon,
                    small,
                    @click="editPopup(index, popup)"
                  )
                    v-icon(color="primary") edit
            v-card-text
              p {{ $t('common.template') }}:
              v-textarea(:value="popup.template", disabled, dark)
      template(#actions="")
        v-btn(
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :loading="submitting",
          :disabled="isDisabled",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ModalWrapper from '../../modal-wrapper.vue';

export default {
  name: MODALS.infoPopupSetting,
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { infoPopups = [] } = this.modal.config;

    return {
      form: {
        popups: cloneDeep(infoPopups),
      },
    };
  },
  methods: {
    addPopup() {
      this.$modals.show({
        name: MODALS.addInfoPopup,
        config: {
          columns: this.config.columns,
          action: popup => this.form.popups.push(popup),
        },
      });
    },

    deletePopup(index) {
      this.$delete(this.form.popups, index);
    },

    editPopup(index, popup) {
      this.$modals.show({
        name: MODALS.addInfoPopup,
        config: {
          columns: this.config.columns,
          popup,
          action: (editedPopup) => {
            this.$set(this.form.popups, index, editedPopup);
          },
        },
      });
    },

    async submit() {
      if (this.config.action) {
        await this.config.action(this.form.popups);
      }

      this.$modals.hide();
    },
  },
};
</script>
