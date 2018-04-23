<template lang='pug'>
  v-container
    v-layout(row wrap)
      v-flex(xs12)
        v-tabs
          v-tab(
            @click='handleTabClick(0)'
          ) {{$t('m_filter_editor.tabs.visual_editor')}}
          v-tab(
            @click='handleTabClick(1)'
          ) {{$t('m_filter_editor.tabs.advanced_editor')}}

    template(v-if='activeTab === 0')
      filter-group(
        initialGroup
        :index = 0
        :condition.sync='filter[0].condition'
        :possibleFields='possibleFields'
        :rules='filter[0].rules'
        :groups='filter[0].groups'
      )

    template(v-if='activeTab === 1')
      v-text-field(
        ref='input'
        v-model='inputValue'
        rows='20'
        :label="$t('m_filter_editor.tabs.advanced_editor')"
        textarea
      )
      v-btn(@click='updateFilter') {{$t('common.parse')}}
</template>

<script>
import { mapState, mapGetters } from 'vuex';
import FilterGroup from '@/components/other/mfilter-editor/filter-group.vue';

export default {
  name: 'mfilter-editor',
  components: {
    FilterGroup,
  },
  data() {
    return {
      newRequest: '',
    };
  },
  computed: {
    ...mapState({
      request: state => state.MFilterEditor.request,
      filter: state => state.MFilterEditor.filter,
      possibleFields: state => state.MFilterEditor.possibleFields,
      activeTab: state => state.MFilterEditor.activeTab,
    }),

    ...mapGetters({
      filter2request: 'MFilterEditor/filter2request',
    }),

    inputValue: {
      get() {
        return JSON.stringify(this.filter2request, undefined, 4);
      },
      set(newVal) {
        this.newRequest = newVal;
      },
    },
  },
  methods: {
    handleTabClick(tab) {
      this.$store.dispatch('MFilterEditor/changeActiveTab', tab);
    },
    updateFilter() {
      if (this.newRequest === '') {
        this.$store.dispatch('MFilterEditor/updateFilter', JSON.parse(JSON.stringify(this.filter2request)));
        return this;
      }
      this.$store.dispatch('MFilterEditor/updateFilter', JSON.parse(this.newRequest));
      return this;
    },
  },
};
</script>
