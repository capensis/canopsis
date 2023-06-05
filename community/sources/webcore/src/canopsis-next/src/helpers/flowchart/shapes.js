import { get, isArray } from 'lodash';

import { LINE_TYPES, SHAPES } from '@/constants';

import uid from '@/helpers/uid';

import { generatePoint } from './points';

/**
 * @typedef {Object} DefaultShapeConnection
 * @property {string} shapeId
 * @property {string} pointId
 * @property {string} side
 */

/**
 * @typedef {Object} DefaultShape
 * @property {string} _id
 * @property {string} type
 * @property {DefaultShapeConnection[]} connections
 * @property {string[]} connectedTo
 * @property {string} text
 * @property {Object} textProperties
 * @property {Object} properties
 * @property {boolean} aspectRatio
 */

/**
 * @typedef {DefaultShape} RectShape
 * @property {number} width
 * @property {number} height
 * @property {number} x
 * @property {number} y
 */

/**
 * @typedef {DefaultShape} LineShape
 * @property {Point[]} points
 * @property {string} lineType
 */

/**
 * @typedef {LineShape} ArrowLineShape
 */

/**
 * @typedef {LineShape} BidirectionalArrowLineShape
 */

/**
 * @typedef {DefaultShape} CircleShape
 * @property {number} diameter
 * @property {number} x
 * @property {number} y
 */

/**
 * @typedef {RectShape} EllipseShape
 */

/**
 * @typedef {RectShape} RhombusShape
 */

/**
 * @typedef {RectShape} ParallelogramShape
 * @property {number} offset
 */

/**
 * @typedef {RectShape} ProcessShape
 * @property {number} offset
 */

/**
 * @typedef {RectShape} DocumentShape
 * @property {number} offset
 */

/**
 * @typedef {RectShape} StorageShape
 * @property {number} radius
 */

/**
 * @typedef {RectShape} ImageShape
 * @property {string} src
 * @property {string} svg
 */

/**
 * @typedef {
 *   RectShape |
 *   CircleShape |
 *   RhombusShape |
 *   EllipseShape |
 *   ParallelogramShape |
 *   ImageShape |
 *   StorageShape |
 *   DocumentShape |
 *   ProcessShape |
 *   LineShape |
 *   ArrowLineShape |
 *   BidirectionalArrowLineShape
 * } Shape
 */

/**
 * Convert default shape to form
 *
 * @param {DefaultShape} shape
 * @returns {DefaultShape}
 */
const defaultShapeToForm = shape => ({
  _id: shape._id ?? uid(),
  connections: shape.connections ?? [],
  connectedTo: shape.connectedTo ?? [],
  text: shape.text ?? '',
  textProperties: {
    ...shape.textProperties,
    fontColor: 'black',
    fontSize: 12,
  },
  properties: shape.properties ?? {},
  aspectRatio: shape.aspectRatio ?? false,
});

/**
 * Convert default rectangle shape to form
 *
 * @param {RectShape} shape
 * @returns {RectShape}
 */
export const rectShapeToForm = shape => ({
  ...defaultShapeToForm(shape),

  type: SHAPES.rect,
  width: shape.width ?? 130,
  height: shape.height ?? 130,
  x: shape.x ?? 0,
  y: shape.y ?? 0,
});

/**
 * Convert default line shape to form
 *
 * @param {LineShape} shape
 * @returns {LineShape}
 */
export const lineShapeToForm = shape => ({
  ...defaultShapeToForm(shape),

  type: SHAPES.line,
  lineType: shape.lineType ?? LINE_TYPES.line,
  points: shape.points ?? [
    generatePoint({
      x: 0,
      y: 100,
    }),
    generatePoint({
      x: 100,
      y: 0,
    }),
  ],
});

/**
 * Convert default arrow line shape to form
 *
 * @param {ArrowLineShape} shape
 * @returns {ArrowLineShape}
 */
export const arrowLineShapeToForm = shape => ({
  ...lineShapeToForm(shape),

  type: SHAPES.arrowLine,
});

/**
 * Convert default bidirectional arrow line shape to form
 *
 * @param {BidirectionalArrowLineShape} shape
 * @returns {BidirectionalArrowLineShape}
 */
export const bidirectionalArrowLineShapeToForm = shape => ({
  ...lineShapeToForm(shape),

  type: SHAPES.bidirectionalArrowLine,
});

/**
 * Convert default circle shape to form
 *
 * @param {CircleShape} shape
 * @returns {CircleShape}
 */
export const circleShapeToForm = shape => ({
  ...defaultShapeToForm(shape),

  type: SHAPES.circle,
  x: shape.x ?? 50,
  y: shape.y ?? 50,
  diameter: shape.diameter ?? 100,
});

/**
 * Convert default ellipse shape to form
 *
 * @param {EllipseShape} shape
 * @returns {EllipseShape}
 */
export const ellipseShapeToForm = shape => ({
  ...rectShapeToForm(shape),

  type: SHAPES.ellipse,
});

/**
 * Convert default rhombus shape to form
 *
 * @param {RhombusShape} shape
 * @returns {RhombusShape}
 */
export const rhombusShapeToForm = shape => ({
  ...rectShapeToForm(shape),

  type: SHAPES.rhombus,
});

/**
 * Convert default parallelogram shape to form
 *
 * @param {ParallelogramShape} shape
 * @returns {ParallelogramShape}
 */
export const parallelogramShapeToForm = shape => ({
  ...rectShapeToForm(shape),

  type: SHAPES.parallelogram,
  offset: shape.offset ?? 50,
});

/**
 * Convert default process shape to form
 *
 * @param {ProcessShape} shape
 * @returns {ProcessShape}
 */
export const processShapeToForm = shape => ({
  ...rectShapeToForm(shape),

  type: SHAPES.process,
  offset: shape.offset ?? 20,
});

/**
 * Convert default document shape to form
 *
 * @param {DocumentShape} shape
 * @returns {DocumentShape}
 */
export const documentShapeToForm = shape => ({
  ...rectShapeToForm(shape),

  type: SHAPES.document,
  offset: shape.offset ?? 20,
});

/**
 * Convert default storage shape to form
 *
 * @param {StorageShape} shape
 * @returns {StorageShape}
 */
export const storageShapeToForm = shape => ({
  ...rectShapeToForm(shape),

  type: SHAPES.storage,
  radius: shape.radius ?? 20,
});

/**
 * Convert default image shape to form
 *
 * @param {ImageShape} shape
 * @returns {ImageShape}
 */
export const imageShapeToForm = shape => ({
  ...rectShapeToForm(shape),

  width: shape.width ?? 40,
  height: shape.height ?? 40,
  type: SHAPES.image,
  src: shape.src,
  svg: shape.svg,
});

/**
 * Convert default shape to form
 *
 * @param {Shape} shape
 * @returns {Shape}
 */
export const shapeToForm = (shape) => {
  const prepare = {
    [SHAPES.rect]: rectShapeToForm,
    [SHAPES.line]: lineShapeToForm,
    [SHAPES.arrowLine]: arrowLineShapeToForm,
    [SHAPES.bidirectionalArrowLine]: bidirectionalArrowLineShapeToForm,
    [SHAPES.circle]: circleShapeToForm,
    [SHAPES.ellipse]: ellipseShapeToForm,
    [SHAPES.rhombus]: rhombusShapeToForm,
    [SHAPES.parallelogram]: parallelogramShapeToForm,
    [SHAPES.process]: processShapeToForm,
    [SHAPES.document]: documentShapeToForm,
    [SHAPES.storage]: storageShapeToForm,
    [SHAPES.image]: imageShapeToForm,
  }[shape.type];

  return prepare(shape);
};

/**
 * Calculate icon position for shape
 *
 * @param {Shape} shape
 * @returns {Point}
 */
export const calculateShapeIconPosition = (shape) => {
  switch (shape.type) {
    case SHAPES.parallelogram:
    case SHAPES.ellipse:
    case SHAPES.process:
    case SHAPES.document:
    case SHAPES.storage:
    case SHAPES.image:
    case SHAPES.rect:
      return {
        x: shape.x + shape.width / 2,
        y: shape.y,
      };
    case SHAPES.rhombus:
      return {
        x: shape.x + shape.width / 2,
        y: shape.y + 5,
      };
    case SHAPES.circle:
      return {
        x: shape.x + shape.diameter / 2,
        y: shape.y,
      };
    default: {
      return {
        x: shape.x,
        y: shape.y,
      };
    }
  }
};

/**
 * Get shape x max and min
 *
 * @param {Shape} shape
 * @returns {{ min: number, max: number }}
 */
export const getShapeXBounds = (shape) => {
  if (shape.points) {
    const xPoints = shape.points.map(({ x }) => x);

    return {
      min: Math.min.apply(null, xPoints),
      max: Math.max.apply(null, xPoints),
    };
  }

  return {
    min: shape.x,
    max: shape.x + (shape.width ?? shape.diameter),
  };
};

/**
 * Get shape y max and min
 *
 * @param {Shape} shape
 * @returns {{ min: number, max: number }}
 */
export const getShapeYBounds = (shape) => {
  if (shape.points) {
    const yPoints = shape.points.map(({ y }) => y);

    return {
      min: Math.min.apply(null, yPoints),
      max: Math.max.apply(null, yPoints),
    };
  }

  return {
    min: shape.y,
    max: shape.y + (shape.width ?? shape.diameter),
  };
};

/**
 * Get shapes max and min coordinate
 *
 * @param {Shape[]} shapes
 * @returns {Object}
 */
export const getShapesBounds = (shapes) => {
  const shapesArray = isArray(shapes) ? shapes : Object.values(shapes);

  return shapesArray.reduce((acc, shape) => {
    const {
      min: minX,
      max: maxX,
    } = getShapeXBounds(shape);
    const {
      min: minY,
      max: maxY,
    } = getShapeYBounds(shape);

    if (minX < acc.min.x) {
      acc.min.x = minX;
    }

    if (minY < acc.min.y) {
      acc.min.y = minY;
    }

    if (maxX > acc.max.x) {
      acc.max.x = maxX;
    }

    if (maxY > acc.max.y) {
      acc.max.y = maxY;
    }

    return acc;
  }, {
    min: {
      x: Infinity,
      y: Infinity,
    },
    max: {
      x: -Infinity,
      y: -Infinity,
    },
  });
};

/**
 * Get property value by shapes array. If all shapes has equal value, will be return value else undefined.
 *
 * @param {Shape[]} shapes
 * @param {string} path
 * @returns {string|number}
 */
export const getPropertyValueByShapes = (shapes, path) => {
  const [firstShape] = shapes;
  const value = get(firstShape, path);

  return shapes?.every(shape => get(shape, path) === value)
    ? value
    : undefined;
};
