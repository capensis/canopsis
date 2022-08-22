import { flowchartPointToForm } from '@/helpers/forms/map';

export const pointsMixin = {
  props: {
    points: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      shownMenu: false,
      pageX: 0,
      pageY: 0,
      clientX: 0,
      clientY: 0,
      addingPoint: undefined,
    };
  },
  computed: {
    isDialogOpened() {
      return !!this.addingPoint;
    },

    contextmenuItems() {
      return [
        {
          text: this.$t('map.addPoint'),
          action: this.openAddPointDialog,
        },
      ];
    },
  },
  methods: {
    emitUpdatePoints(data) {
      this.$emit('update:points', data);
    },

    updatePointInModel(data) {
      this.emitUpdatePoints(
        this.points.map(point => (point._id === data._id ? data : point)),
      );
    },

    addPointToModel(data) {
      this.emitUpdatePoints([...this.points, data]);
    },

    removePointFromModel(data) {
      this.emitUpdatePoints(
        this.points.filter(point => data._id !== point._id),
      );
    },

    openContextmenu() {
      this.shownMenu = true;
    },

    closeContextmenu() {
      this.shownMenu = false;
    },

    setPageOffsetByEvent(event) {
      this.pageX = event.pageX;
      this.pageY = event.pageY - window.scrollY;
      this.clientX = event.clientX;
      this.clientY = event.clientY;
    },

    handleContextmenu(event) {
      if (this.shownMenu || this.isDialogOpened) {
        return;
      }

      this.setPageOffsetByEvent(event);
      this.openContextmenu();
    },

    openAddPointDialog() {
      const { x, y } = this.normalizeCursor({ x: this.clientX, y: this.clientY });

      this.addingPoint = flowchartPointToForm({
        x,
        y,
      });

      this.closeContextmenu();
    },

    closePointDialog() {
      this.addingPoint = undefined;
      this.closeContextmenu();
    },

    submitPointDialog(data) {
      if (this.editingPoint) {
        this.updatePointInModel(data);
      } else {
        this.addPointToModel(data);
      }

      this.closePointDialog();
    },
  },
};
