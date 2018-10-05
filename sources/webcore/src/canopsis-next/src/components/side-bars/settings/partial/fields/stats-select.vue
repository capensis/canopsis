<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ $t('settings.statsSelect') }}
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
            v-list-group.my-2(v-if="options.length > 0")
              v-list-tile(slot="activator") {{ $t('common.options') }}
              template(v-for="option in options")
                v-switch(
                v-show="option === 'recursive'"
                :label="$t('common.recursive')",
                v-model="form.parameters.recursive",
                hide-details
                )
                v-select(
                v-show="option === 'states'"
                :placeholder="$t('common.states')",
                :items="stateTypes",
                v-model="form.parameters.states",
                multiple,
                chips,
                hide-details
                )
                v-combobox(
                v-show="option === 'authors'"
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
            v-alert.my-1(:value="error", type="error") {{ error }}
            v-btn.ma-0(@click="submitClick") {{ editing ? $t('common.edit') : $t('common.add') }}
            v-btn.ma-0(v-if="editing", @click="stopEditing") {{ $t('common.quitEditing') }}

      v-container
        v-list(dark)
          v-list-group.my-1.grey(v-for="(stat, key) in value", :key="key")
            v-list-tile(slot="activator") {{ key }}
              v-layout(justify-end)
                v-btn.green.darken-4.white--text(@click.stop="e => editStat(key)", fab, small, depressed)
                  v-icon edit
                v-btn.red.darken-4.white--text(@click.stop="e => deleteStat(key)", fab, small, depressed)
                  v-icon delete
            v-container(fluid)
              p {{ $t('common.stat') }}: {{ stat.stat }}
              p {{ $t('common.trend') }}: {{ stat.trend }}
              p {{ $t('common.parameters') }}: {{ stat.parameters }}

</template>

<script>
import omit from 'lodash/omit';
import set from 'lodash/set';
import unset from 'lodash/unset';
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
      editing: false,
      editingStatTitle: '',
      form: {
        stat: STATS_TYPES.alarmsCreated,
        title: '',
        trend: true,
        parameters: {
          sla: '',
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
    options() {
      if (this.form.stat) {
        return this.form.stat.options;
      }
      return [];
    },
  },
  methods: {
    async addStat() {
      if (this.value[this.form.title]) {
        this.error = this.$t('settings.statSelector.error.alreadyExist');
      } else {
        const isFormValid = await this.$validator.validateAll();

        if (isFormValid) {
          const newValue = { ...this.value };
          this.error = '';
          const newStat = omit(this.form, ['title', 'parameters', 'stat']);
          newStat.stat = this.form.stat.value;
          newStat.parameters = {};
          this.options.forEach((option) => {
            newStat.parameters[option] = this.form.parameters[option];
          });
          this.$emit('input', set(newValue, this.form.title, newStat));
        }
      }
    },

    deleteStat(stat) {
      const newValue = { ...this.value };
      unset(newValue, stat);
      this.$emit('input', newValue);
    },

    editStat(key) {
      this.editing = true;
      this.editingStatTitle = key;
      this.form = { ...this.value[key], title: key };
    },

    stopEditing() {
      this.editing = false;
      this.editingStatTitle = '';
    },

    submitClick() {
      if (this.editing) {
        // Delete the stat that we want to edit
        const newValue = { ...this.value };
        const { editingStatTitle } = { ...this };
        unset(newValue, editingStatTitle);
        // Set the edited stat in newValue object, and send it to parent with input event
        this.$emit('input', set(newValue, this.form.title, omit(this.form, ['title'])));
      } else {
        this.addStat();
      }
    },
  },
};
</script>

