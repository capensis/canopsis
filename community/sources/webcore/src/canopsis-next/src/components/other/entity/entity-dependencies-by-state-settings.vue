<template>
  <div class="entity-dependencies-by-state-settings">
    <network-graph
      ref="networkGraph"
      :options="options"
      :node-html-label-options="nodeHtmlLabelsOptions"
      class="fill-height black--text"
      ctrl-wheel-zoom
    />
  </div>
</template>

<script>
import { omit } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { ROOT_CAUSE_DIAGRAM_OPTIONS, ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS } from '@/constants';

import { normalizeTreeOfDependenciesMapEntities } from '@/helpers/entities/map/list';
import { isEntityEventsStateSettings } from '@/helpers/entities/entity/entity';
import { convertSortToRequest } from '@/helpers/entities/shared/query';
import { getButtonHTML, getEntityNodeElementHTML } from '@/helpers/entities/entity/cytoscape';

import { entitiesEntityDependenciesMixin } from '@/mixins/entities/entity-dependencies';

import NetworkGraph from '@/components/common/chart/network-graph.vue';

export default {
  components: { NetworkGraph },
  mixins: [entitiesEntityDependenciesMixin],
  props: {
    entity: {
      type: Object,
      required: true,
    },
    colorIndicator: {
      type: String,
      required: false,
    },
    impact: {
      type: Boolean,
      required: false,
    },
    columns: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      metaByEntityId: {},
      entitiesById: normalizeTreeOfDependenciesMapEntities([{ entity: this.entity, pinned_entities: [] }]),
    };
  },
  computed: {
    isEventsStateSettings() {
      return isEntityEventsStateSettings(this.entity);
    },

    entitiesWithDependencies() {
      return Object.values(this.entitiesById).filter(entity => entity.dependencies);
    },

    entitiesElements() {
      const rootElement = this.entitiesById[this.entity._id];
      const { entity, dependencies = [] } = rootElement;

      const elements = [
        {
          group: 'nodes',
          data: {
            id: entity._id,
            entity,
            root: true,
            opened: true,
          },
        },
      ];

      if (isEntityEventsStateSettings(entity)) {
        elements.push(...this.getEventsNodeElementByEntity(entity));

        return elements;
      }

      elements.push(...this.getEntityDependenciesElement(entity, dependencies, [entity._id]));

      return elements;
    },

    styleOption() {
      return [
        {
          selector: 'node',
          style: {
            width: ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize,
            height: ROOT_CAUSE_DIAGRAM_OPTIONS.nodeSize,
          },
        },
        {
          selector: 'node[showMore]',
          style: {
            'background-opacity': 0,
            'border-width': 0,
            width: 128,
            height: 34,
          },
        },
        {
          selector: 'node[isEvents]',
          style: {
            width: 30,
            height: 30,
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

    nodeHtmlLabelsOptions() {
      const getShowMoreContent = (node) => {
        const { entity } = node;
        const meta = this.metaByEntityId[entity._id] ?? {};

        const fetchedEntities = meta.page * meta.per_page;

        /**
         * TODO: Should be replaced on translation
         */
        return getButtonHTML(`Show more (${fetchedEntities} of ${meta.total_count})`);
      };

      return [
        {
          query: 'node',
          valign: 'center',
          halign: 'center',
          tpl: getEntityNodeElementHTML,
        },
        {
          query: 'node[showMore]',
          valign: 'center',
          halign: 'center',
          tpl: getShowMoreContent,
        },
      ];
    },

    options() {
      const options = {
        ...omit(ROOT_CAUSE_DIAGRAM_OPTIONS, ['nodeSize']),

        style: this.styleOption,
        elements: this.entitiesElements,
      };

      if (this.entitiesElements.length) {
        options.layout = {
          ...ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
        };
      }

      return options;
    },
  },
  watch: {
    entity() {
      this.entitiesById = normalizeTreeOfDependenciesMapEntities([{ entity: this.entity, pinned_entities: [] }]);

      /**
       * TODO: investigate this behavior in the future
       */
      setTimeout(() => this.resetLayout(), 1000);
    },
  },
  async mounted() {
    this.$refs.networkGraph.$cy.on('tap', this.tapHandler);

    if (!this.isEventsStateSettings) {
      await this.fetchDependencies(this.entity._id);
    }

    this.runLayout();
  },
  beforeDestroy() {
    this.$refs.networkGraph.$cy.off('tap', this.tapHandler);
  },
  methods: {
    getEventsNodeElementByEntity(entity) {
      const eventsNodeId = `${entity._id}_events-node`;

      return [
        {
          group: 'nodes',
          data: {
            entity,
            id: eventsNodeId,
            isEvents: true,
          },
        },
        {
          group: 'edges',
          data: {
            source: entity._id,
            target: eventsNodeId,
          },
        },
      ];
    },

    getEntityDependenciesElement(entity, dependenciesIds = [], handledDependenciesIds = []) {
      const dependenciesNodes = dependenciesIds.reduce((acc, childId) => {
        const { dependencies: childDependenciesIds = [], entity: child } = this.entitiesById[childId];

        const isCycle = handledDependenciesIds.includes(childId);

        const hasDependencies = !!childDependenciesIds.length;

        if (!isCycle) {
          const childDependencies = this.getEntityDependenciesElement(
            child,
            childDependenciesIds,
            [...handledDependenciesIds, childId],
          );

          acc.push(
            {
              group: 'nodes',
              data: {
                id: childId,
                entity: child,
                root: hasDependencies,
                opened: hasDependencies,
              },
            },
            ...childDependencies,
          );
        }

        acc.push(
          {
            group: 'edges',
            data: {
              source: entity._id,
              target: childId,
            },
          },
        );

        return acc;
      }, []);

      dependenciesNodes.push(
        ...this.getShowMoreElements(entity),
      );

      return dependenciesNodes;
    },

    getShowMoreElements(entity) {
      const meta = this.metaByEntityId[entity._id];

      if (!meta || meta.page >= meta.page_count) {
        return [];
      }

      const showMoreId = `show-all-${entity._id}`;

      return [
        {
          group: 'nodes',
          data: {
            id: showMoreId,
            entity,
            showMore: true,
          },
        },
        {
          group: 'edges',
          data: {
            id: `show-all-edge-${entity._id}`,
            source: entity._id,
            target: showMoreId,
          },
        },
      ];
    },

    /**
     * Remove old elements and add new elements to network graph
     */
    resetLayout() {
      this.$refs.networkGraph.$cy.elements().remove();
      this.$refs.networkGraph.$cy.add(this.entitiesElements);
      this.runLayout();
    },

    /**
     * Run 'cise' layout for rerender clusters
     */
    async runLayout() {
      if (this.$refs.networkGraph.$cy.nodes().empty()) {
        return;
      }

      try {
        this.$refs.networkGraph.$cy.layout({
          ...ROOT_CAUSE_DIAGRAM_LAYOUT_OPTIONS,
        }).run();
      } catch (err) {
        console.warn(err);
      }
    },

    /**
     * Show dependencies for node
     *
     * @param {Object} target
     * @returns {Promise<void>}
     */
    async showDependencies(target) {
      const { id } = target.data();
      const { page } = this.metaByEntityId[id] ?? {};
      const newPage = page ? page + 1 : 1;

      target.data({
        pending: true,
      });

      const { data, meta } = await this.fetchDependenciesList({
        id,
        params: {
          page: newPage,
          limit: PAGINATION_LIMIT,
          with_flags: true,
          /**
           * TODO: Api doesn't support multi sort
           */
          ...convertSortToRequest(['last_update_date', 'state']),
        },
      });

      target.data({
        pending: false,
      });

      this.$set(this.metaByEntityId, id, meta);

      const ids = data.map((item) => {
        let newEntityItem = { entity: item };

        if (this.entitiesById[item._id]) {
          newEntityItem = {
            ...this.entitiesById[item._id],

            entity: {
              ...newEntityItem,
              ...this.entitiesById[item._id].entity,
            },
          };
        }

        this.$set(this.entitiesById, item._id, newEntityItem);

        return item._id;
      });

      const previousDeps = this.entitiesById[id].dependencies ?? [];

      this.$set(this.entitiesById[id], 'dependencies', [
        ...previousDeps,
        ...ids,
      ]);

      this.$refs.networkGraph.$cy.elements('*').remove();
      this.$refs.networkGraph.$cy.add(this.entitiesElements);
    },

    /**
     * Method for dependencies fetching for special node
     *
     * @param {string} id
     */
    async fetchDependencies(id) {
      const target = this.$refs.networkGraph.$cy.getElementById(id);
      const { pending } = target.data();

      if (pending) {
        return;
      }

      await this.showDependencies(target);

      this.runLayout();
    },

    /**
     * Handler for tap event on whole cytoscape canvas
     *
     * @param {Object} target
     * @param {MouseEvent} originalEvent
     */
    tapHandler({ target, originalEvent }) {
      const { entity, showMore, cycle } = target.data();

      if (cycle) {
        return;
      }

      if (originalEvent.target.classList.contains('v-badge__badge')) {
        const { id } = originalEvent.target.dataset;

        if (id) {
          this.fetchDependencies(id);

          return;
        }
      }

      if (!showMore || !entity) {
        return;
      }

      this.fetchDependencies(entity._id);
    },
  },
};
</script>

<style lang="scss">
.entity-dependencies-by-state-settings {
  position: relative;
  height: 800px;
  width: 100%;
  border-radius: 5px;
  background: white;

  &__node-progress {
    position: absolute;
    inset: 0;
  }

  &__fetch-dependencies {
    width: 100%;
    height: 100%;
    border-radius: 50%;
  }

  canvas[data-id='layer0-selectbox'] { // Hide selectbox layer from cytoscape
    display: none;
  }

  .v-badge__badge {
    top: -7px;
    right: -7px;

    * {
      pointer-events: none;
    }
  }
}
</style>
