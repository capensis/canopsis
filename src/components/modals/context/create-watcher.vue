<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t(config.title) }}
    v-form
      v-layout(wrap, justify-center)
        v-flex(xs11)
          v-text-field(
            :label="$t('modals.createWatcher.displayName')",
            v-model="form.name",
            data-vv-name="name",
            v-validate="'required'",
            :error-messages="errors.collect('name')",
          )
      v-layout(wrap, justify-center)
        v-flex(xs11)
          h3.text-xs-center {{ $t('filterEditor.title') }}
          v-divider
          filter-editor(v-model="form.mfilter")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit", :loading="submitting", :disabled="submitting") {{ $t('common.submit') }}
</template>

<script>
import uid from '@/helpers/uid';
import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import popupMixin from '@/mixins/popup';
import { MODALS, ENTITIES_TYPES } from '@/constants';

export default {
  name: MODALS.createWatcher,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FilterEditor,
  },
  mixins: [modalInnerMixin, entitiesContextEntityMixin, popupMixin],
  data() {
    const { item } = this.modal.config;

    let form = {
      name: '',
      mfilter: '{}',
    };

    if (item) {
      form = { ...item };
    }

    return {
      form,
      submitting: false,
    };
  },
  methods: {
    async submit() {
      this.submitting = true;
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          ...this.form,
          _id: this.config.item && !this.config.isDuplicating ? this.config.item._id : uid(),
          display_name: this.form.name,
          type: ENTITIES_TYPES.watcher,
        };

        try {
          await this.config.action(data);
          this.refreshContextEntitiesLists();
          this.hideModal();

          if (this.config.item && !this.config.isDuplicating) {
            this.addSuccessPopup({ text: this.$t('modals.createWatcher.success.edit') });
          } else if (this.config.isDuplicating) {
            this.addSuccessPopup({ text: this.$t('modals.createWatcher.success.duplicate') });
          } else {
            this.addSuccessPopup({ text: this.$t('modals.createWatcher.success.create') });
          }
        } catch (err) {
          this.addErrorPopup({ text: this.$t('error.default') });
          console.error(err);
        }

        this.submitting = false;
      }
    },
  },
};
</script>
