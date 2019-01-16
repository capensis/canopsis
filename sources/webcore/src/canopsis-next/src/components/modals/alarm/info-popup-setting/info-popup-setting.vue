<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.infoPopupSetting.title') }}
    v-card-text
      v-layout(justify-end)
        v-btn(@click="addPopup", icon, fab, small, color="secondary")
          v-icon add
      v-layout(column)
        v-card.my-1(v-for="(popup, index) in popups", :key="index", flat, color="secondary white--text")
          v-card-title
            v-layout(justify-space-between)
              div {{ $t('modals.infoPopupSetting.column') }}: {{ popup.column }}
              div
                v-btn(@click="deletePopup(index)", icon, small)
                  v-icon(color="error") delete
                v-btn(@click="editPopup(index, popup)", icon, small)
                  v-icon(color="primary") edit
          v-card-text
            p {{ $t('modals.infoPopupSetting.template') }}:
            v-textarea(:value="popup.template", :disabled="true", dark)
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit", type="submit") {{ $t('common.submit') }}
</template>

<script>
import pullAt from 'lodash/pullAt';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import TextEditor from '@/components/other/text-editor/text-editor.vue';

export default {
  name: MODALS.infoPopupSetting,
  components: {
    TextEditor,
  },
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
      const popups = [...this.popups];
      pullAt(popups, index);

      this.popups = popups;
    },

    editPopup(index, popup) {
      this.showModal({
        name: MODALS.addInfoPopup,
        config: {
          columns: this.config.columns,
          popup,
          action: (editedPopup) => {
            const popups = [...this.popups];
            popups[index] = editedPopup;

            this.popups = popups;
          },
        },
      });
    },

    submit() {
      if (this.config.action) {
        this.config.action(this.popups);
      }

      this.hideModal();
    },
  },
};
</script>
