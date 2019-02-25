<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
    v-form
      v-container
        v-card.mb-2
          v-container.pt-0(fluid)
            v-select(
            v-model="form.stat",
            hide-details,
            :items="statsTypes",
            return-object,
            )
            v-text-field(
            :placeholder="$t('common.title')",
            v-model="form.title",
            :error-messages="errors.collect('title')",
            v-validate="'required'",
            data-vv-name="title",
            )
            v-switch(
            :label="$t('common.trend')",
            v-model="form.trend",
            hide-details
            )
            v-list
              v-list-group.my-2
                v-list-tile(slot="activator") {{ $t('common.options') }}
                template(v-for="option in options")
                  v-switch(
                  v-show="option === 'recursive'",
                  :label="$t('common.recursive')",
                  v-model="form.parameters.recursive",
                  hide-details
                  )
                  v-select(
                  v-show="option === 'states'",
                  :placeholder="$t('common.states')",
                  :items="stateTypes",
                  v-model="form.parameters.states",
                  multiple,
                  chips,
                  hide-details
                  )
                  v-combobox(
                  v-show="option === 'authors'",
                  :placeholder="$t('common.authors')",
                  v-model="form.parameters.authors",
                  hide-details,
                  chips,
                  multiple
                  )
                  v-text-field(
                  v-show="option === 'sla'",
                  :placeholder="$t('common.sla')",
                  v-model="form.parameters.sla",
                  hide-details
                  )
      v-divider
      v-layout.py-1(justify-end)
        v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, STATS_TYPES, ENTITIES_STATES } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

export default {
  name: MODALS.addStat,
  $_veeValidate: {
    validator: 'new',
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      form: {
        stat: STATS_TYPES.alarmsCreated,
        title: '',
        trend: true,
        parameters: {},
      },
      error: '',
    };
  },
  computed: {
    /**
     * Get stats different types from constant, and return an object with stat's value and stat's translated title
     */
    statsTypes() {
      return Object.values(STATS_TYPES)
        .map(item => ({ value: item.value, text: this.$t(`stats.types.${item.value}`), options: item.options }));
    },
    stateTypes() {
      return Object.keys(ENTITIES_STATES)
        .map(item => ({ value: ENTITIES_STATES[item], text: item }));
    },
    options() {
      if (this.form.stat) {
        return this.form.stat.options;
      }

      return [];
    },
  },
  mounted() {
    if (this.config.stat) {
      this.form = { ...this.config.stat, title: this.config.statTitle };
    }
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid && this.config.action) {
        await this.config.action(this.form);
        this.hideModal();
      }
    },
  },
};
</script>
