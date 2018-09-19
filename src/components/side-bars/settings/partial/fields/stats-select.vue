<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsSelect') }}
    v-form(ref="newStatForm")
      v-container
        v-card.mb-2
          v-container.pt-0(fluid)
            v-select(
            v-model="form.stat",
            hide-details,
            :items="statsTypes",
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
            v-list-group.my-2
              v-list-tile(slot="activator") Options
              v-switch(
              :label="$t('common.recursive')",
              v-model="form.parameters.recursive",
              hide-details
              )
              v-select(
              :placeholder="$t('common.states')",
              :items="stateTypes",
              v-model="form.parameters.states",
              multiple,
              chips,
              hide-details
              )
              v-combobox(
              :placeholder="$t('common.authors')",
              v-model="form.parameters.authors",
              hide-details,
              chips,
              multiple
              )
              v-text-field(:placeholder="$t('common.sla')", v-model="form.sla", hide-details)
            v-btn.ma-0(@click="addStat") Add stat

      v-list-group(v-for="(stat, key) in value", :key="key")
        v-list-tile(slot="activator") {{ key }}
        v-container(fluid) {{ stat }}

</template>

<script>
import omit from 'lodash/omit';
import set from 'lodash/set';
import { STATS_TYPES, ENTITIES_STATES } from '@/constants';

export default {
  $_veeValidate: {
    validator: 'new',
  },
  props: {
    value: {
      type: Object,
    },
  },
  data() {
    return {
      form: {
        stat: 'alarms_created',
        title: '',
        sla: '',
        trend: true,
        parameters: {
          states: [],
          recursive: false,
          authors: [],
        },
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
      return Object.keys(ENTITIES_STATES).map(item => ({ value: ENTITIES_STATES[item], text: item }));
    },
  },
  methods: {
    async addStat() {
      if (this.value[this.form.title]) {
        this.error = 'Stat with this title already exists';
      } else {
        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          this.error = '';
          this.$emit('input', set(this.value, this.form.title, omit(this.form, ['title'])));
          this.$refs.newStatForm.reset();
        }
      }
    },
  },
};
</script>

