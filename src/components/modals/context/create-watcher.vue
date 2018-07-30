<template lang="pug">
  v-card
    v-card-title
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
          h3.text-xs-center {{ $t('mFilterEditor.title') }}
          v-divider
          filter-editor
      v-layout(wrap)
        v-flex(xs12)
          entities-select(label="Impacts", :entities="form.impact", @updateEntities="updateImpact")
        v-flex(xs12)
          entities-select(label="Dependencies", :entities="form.depends", @updateEntities="updateDepends")
      v-layout
        v-flex(xs3)
          v-btn.green.darken-4.white--text(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import FilterEditor from '@/components/other/filter-editor/filter-editor.vue';
import EntitiesSelect from '@/components/modals/context/entities-select.vue';
import modalInnerMixin from '@/mixins/modal/modal-inner';
import { ENTITIES_TYPES } from '@/constants';

const { mapActions: watcherMapActions } = createNamespacedHelpers('watcher');
const { mapGetters: filterEditorMapGetters } = createNamespacedHelpers('mFilterEditor');

export default {
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FilterEditor,
    EntitiesSelect,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        name: '',
        mfilter: '',
        impact: [],
        depends: [],
      },
    };
  },
  computed: {
    ...filterEditorMapGetters(['request']),
  },
  mounted() {
    if (this.config.item) {
      this.form = { ...this.config.item };
    }
  },
  methods: {
    ...watcherMapActions(['create', 'edit']),
    updateImpact(entities) {
      this.form.impact = entities;
    },
    updateDepends(entities) {
      this.form.depends = entities;
    },
    async submit() {
      const formIsValid = await this.$validator.validateAll();

      if (formIsValid) {
        // If there's an item, means we're editing. If there's not, we're creating an entity
        let formData = { ...this.form };
        if (this.config.item) {
          formData = {
            _id: this.config.item._id,
            display_name: this.form.name,
            type: ENTITIES_TYPES.watcher,
            impact: this.form.impact,
            depends: this.form.depends,
            mfilter: JSON.stringify(this.request),
          };
          await this.edit({ watcher_id: formData._id, data: formData });
        } else {
          formData = {
            ...this.form,
            _id: this.form.name,
            display_name: this.form.name,
            type: ENTITIES_TYPES.watcher,
            impact: this.form.impact,
            depends: this.form.depends,
            mfilter: JSON.stringify(this.request),
          };
          await this.create(formData);
        }
        this.hideModal();
      }
    },
  },
};
</script>
