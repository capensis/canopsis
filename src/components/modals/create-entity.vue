<template lang="pug">
  v-card
    v-card-title.blue.darken-4.white--text.text-xs-center
      h2 {{ $t(config.title) }}
    v-card-text
      v-container
        v-layout(row)
          v-text-field(
            :label="$t('common.name')",
            :error-messages="errors.collect('name')",
            v-model="form.name",
            v-validate="'required'",
            data-vv-name="name",
          )
        v-layout.mt-2(row)
          v-text-field(
            :label="$t('common.description')",
            :error-messages="errors.collect('description')",
            v-model="form.description",
            v-validate="'required'",
            data-vv-name="description",
            multi-line,
          )
        v-layout.mt-2(row, align-center)
          v-switch(:label="$t('common.enabled')", v-model="form.enabled")
          v-select.pa-0(
            :items="types",
            v-model="form.type",
            label="Type",
            single-line,
          )
    entities-select(label="Impacts", :entities.sync=entities)
    entities-select(label="Dependencies", :entities.sync=entities)
    v-card-actions
      v-btn(@click.prevent="submit", color="blue darken-4 white--text") {{ $t('common.submit') }}
      v-btn(
        @click.prevent="manageInfos",
        color="blue darken-4 white--text"
      ) {{ $t('modals.createEntity.fields.manageInfos') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import EntitiesSelect from '@/components/other/context/actions/create-entities/entities-select.vue';
import ModalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

const { mapActions: entitiesMapActions } = createNamespacedHelpers('context');

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
  },
  components: { EntitiesSelect },
  mixins: [ModalInnerMixin],
  data() {
    return {
      types: [
        {
          text: this.$t('modals.createEntity.fields.types.connector'),
          value: 'connector',
        },
        {
          text: this.$t('modals.createEntity.fields.types.component'),
          value: 'component',
        },
        {
          text: this.$t('modals.createEntity.fields.types.resource'),
          value: 'resource',
        },
      ],
      showValidationErrors: true,
      enabled: true,
      form: {
        name: '',
        description: '',
        type: '',
        enabled: true,
        depends: [],
        impact: [],
        infos: [],
      },
    };
  },
  mounted() {
    if (this.config.item) {
      this.form = { ...this.config.item.props };
    }
  },
  methods: {
    ...entitiesMapActions({
      createEntity: 'create',
      editEntity: 'edit',
    }),
    updateImpact(entities) {
      this.form.impacts = entities.map(entity => entity._id);
    },
    updateDependencies(entities) {
      this.form.dependencies = entities.map(entity => entity._id);
    },
    async submit() {
      const formIsValid = await this.$validator.validateAll();
      if (formIsValid) {
        // If there's an item, means we're editing. If there's not, we're creating an entity
        if (this.config.item) {
          this.editEntity({ data: this.form });
        } else {
          const formData = { ...this.form, _id: this.form.name };
          this.createEntity({ data: formData });
        }
        this.hideModal();
      }
    },

  },
};
</script>

<style scoped>
  .tooltip {
    flex: 1 1 auto;
  }
  .impact {
    background-color: grey;
  }
</style>
