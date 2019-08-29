<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
    v-tabs(slider-color="primary")
      v-tab(
        v-for="tab in tabs",
        :key="tab.name"
      ) {{ tab.name }}
      v-tab-item
        v-form
          v-layout(wrap, justify-center)
            v-flex(xs11)
              v-text-field(
                :label="$t('modals.createWatcher.displayName')",
                v-model="form.name",
                :error-messages="errors.collect('name')",
                data-vv-name="name",
                v-validate="'required'"
              )
          v-layout(wrap, justify-center)
            v-flex(xs11)
              template(v-if="stack === $constants.CANOPSIS_STACK.go")
                v-textarea(
                  label="Output template",
                  v-model="form.output_template",
                  :error-messages="errors.collect('output_template')",
                  data-vv-name="output_template",
                  v-validate="'required'"
                )
      v-tab-item
        v-card
          v-card-text
            patterns-list(v-if="stack === $constants.CANOPSIS_STACK.go", v-model="form.entities")
            filter-editor(v-else, v-model="form.mfilter", required)
      v-tab-item
        manage-infos(v-model="form.infos")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit", :loading="submitting", :disabled="submitting") {{ $t('common.submit') }}
</template>

<script>
import { omit } from 'lodash';

import { MODALS, ENTITIES_TYPES, CANOPSIS_STACK } from '@/constants';

import uuid from '@/helpers/uuid';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import entitiesInfoMixin from '@/mixins/entities/info';

import FilterEditor from '@/components/other/filter/editor/filter-editor.vue';
import PatternsList from '@/components/other/shared/patterns-list/patterns-list.vue';

import ManageInfos from './partial/manage-infos.vue';

export default {
  name: MODALS.createWatcher,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    FilterEditor,
    PatternsList,
    ManageInfos,
  },
  mixins: [modalInnerMixin, entitiesContextEntityMixin, entitiesInfoMixin],
  data() {
    const { item } = this.modal.config;

    let form = {
      name: '',
      mfilter: '{}',
      infos: {},
      impact: [],
      depends: [],
      entities: [],
      output_template: '',
    };

    if (item) {
      form = { ...item };
    }

    return {
      form,
      submitting: false,
    };
  },
  computed: {
    tabs() {
      return [
        { name: this.$t('modals.createEntity.fields.form') },
        { name: this.stack === CANOPSIS_STACK.go ? this.$t('eventFilter.pattern') : this.$t('common.filter') },
        { name: this.$t('modals.createEntity.fields.manageInfos') },
      ];
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        this.submitting = true;
        let data = {};

        if (this.stack === CANOPSIS_STACK.go) {
          data = {
            ...omit(this.form, ['mfilter', 'impact', 'depends']),
            _id: this.config.item && !this.config.isDuplicating ? this.config.item._id : uuid('watcher'),
            name: this.form.name,
            type: ENTITIES_TYPES.watcher,
            state: {
              method: 'worst',
            },
          };
        } else {
          data = {
            ...omit(this.form, ['entities', 'output_template']),
            _id: this.config.item && !this.config.isDuplicating ? this.config.item._id : uuid('watcher'),
            infos: this.form.infos,
            display_name: this.form.name,
            type: ENTITIES_TYPES.watcher,
          };
        }

        try {
          await this.config.action(data);
          this.refreshContextEntitiesLists();

          this.hideModal();
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

<style>

</style>
