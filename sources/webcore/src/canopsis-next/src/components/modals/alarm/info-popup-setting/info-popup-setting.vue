<template lang="pug">
  v-card(data-test="infoPopupSettingModal")
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.infoPopupSetting.title') }}
    v-card-text
      v-layout(justify-end)
        v-btn(
          data-test="infoPopupAddPopup",
          icon,
          fab,
          small,
          color="secondary",
          @click="addPopup"
        )
          v-icon add
      v-layout(column)
        v-card.my-1(v-for="(popup, index) in popups", :key="index", flat, color="secondary white--text")
          v-card-title
            v-layout(justify-space-between)
              div {{ $t('modals.infoPopupSetting.column') }}: {{ popup.column }}
              div
                v-btn(
                  data-test="infoPopupDeletePopup",
                  icon,
                  small,
                  @click="deletePopup(index)"
                )
                  v-icon(color="error") delete
                v-btn(
                  data-test="infoPopupEditPopup",
                  icon,
                  small,
                  @click="editPopup(index, popup)"
                )
                  v-icon(color="primary") edit
          v-card-text
            p {{ $t('modals.infoPopupSetting.template') }}:
            v-textarea(:value="popup.template", disabled, dark)
    v-divider
    v-layout.py-1(justify-end)
      v-btn(
        data-test="infoPopupCancelButton",
        depressed,
        flat,
        @click="hideModal"
      ) {{ $t('common.cancel') }}
      v-btn.primary(
        data-test="infoPopupSubmitButton",
        type="submit",
        @click="submit"
      ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.infoPopupSetting,
  mixins: [modalInnerMixin],
  data() {
    return {
      popups: [],
    };
  },
  mounted() {
    if (this.config) {
      this.popups = [...this.config.infoPopups];
    }
  },
  methods: {
    addPopup() {
      this.showModal({
        name: MODALS.addInfoPopup,
        config: {
          columns: this.config.columns,
          action: popup => this.popups.push(popup),
        },
      });
    },

    deletePopup(index) {
      this.$delete(this.popups, index);
    },

    editPopup(index, popup) {
      this.showModal({
        name: MODALS.addInfoPopup,
        config: {
          columns: this.config.columns,
          popup,
          action: (editedPopup) => {
            this.$set(this.popups, index, editedPopup);
          },
        },
      });
    },

    async submit() {
      if (this.config.action) {
        await this.config.action(this.popups);
      }

      this.hideModal();
    },
  },
};
</script>
