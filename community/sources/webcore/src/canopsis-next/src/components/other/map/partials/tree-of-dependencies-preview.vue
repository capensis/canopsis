<template lang="pug">
  c-zoom-overlay
    network-graph.fill-height(
      ref="networkGraph",
      :options="options",
      :node-html-label-options="nodeHtmlLabelsOptions",
      ctrl-wheel-zoom
    )
</template>

<script>
import { keyBy } from 'lodash';

import { COLORS } from '@/config';
import {
  TREE_OF_DEPENDENCIES_GRAPH_OPTIONS,
  TREE_OF_DEPENDENCIES_GRAPH_LAYOUT_OPTIONS,
} from '@/constants';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

export default {
  components: { NetworkGraph },
  props: {
    entities: {
      type: Array,
      required: true,
    },
    impact: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    rootEntitiesById() {
      return keyBy(this.entities, 'data._id');
    },

    cytoscapeClusters() {
      return this.entities.map(({ data = {}, pinned = [] } = {}) => [
        data._id,

        ...pinned.filter(({ _id: id }) => !this.rootEntitiesById[id]).map(({ _id: id }) => id),
      ]);
    },

    styleOption() {
      return [
        {
          selector: 'node',
          style: {
            'font-size': '10px',
            'background-color': COLORS.secondary,
            'border-color': COLORS.secondary,
          },
        },
        {
          selector: 'edge',
          style: {
            width: 2,
            'curve-style': 'bezier',
            'line-color': 'silver',
          },
        },
      ];
    },

    countProperty() {
      return this.impact ? 'impacts_count' : 'depends_count';
    },

    nodeHtmlLabelsOptions() {
      const getIcon = ({ entity, opened }) => {
        const count = entity[this.countProperty];

        return count
          ? `<i class="v-icon material-icons theme--light white--text">${opened ? 'remove' : 'add'}</i>`
          : '';
      };

      return [
        {
          query: 'node',
          valign: 'center',
          halign: 'center',
          tpl: data => `<div class="secondary v-btn__content" style="width: 30px; height: 30px; border-radius: 50%;">${getIcon(data)}</div>`,
        },
      ];
    },

    entitiesElements() {
      return this.entities.reduce((acc, { data, pinned = [] }) => {
        acc.push(
          {
            group: 'nodes',
            data: { id: data._id, entity: data },
          },
        );

        pinned.forEach((entity) => {
          acc.push(
            {
              group: 'nodes',
              data: { id: entity._id, entity },
            },
            {
              group: 'edges',
              data: { source: data._id, target: entity._id },
            },
          );
        });

        return acc;
      }, []);
    },

    options() {
      return {
        ...TREE_OF_DEPENDENCIES_GRAPH_OPTIONS,

        layout: {
          ...TREE_OF_DEPENDENCIES_GRAPH_LAYOUT_OPTIONS,

          clusters: this.cytoscapeClusters,
        },

        style: this.styleOption,
        elements: this.entitiesElements,
      };
    },

    /*    nodeHtmlLabelsOptions() {
      const tpl = ({ entity }) => `<div>${entity.type === ENTITY_TYPES.service ? entity.name : entity._id}</div>`;

      return [
        {
          tpl,
          query: 'node',
          valign: 'top',
          halign: 'center',
          valignBox: 'top',
          halignBox: 'center',
        },
      ];
    }, */
  },
};
</script>
