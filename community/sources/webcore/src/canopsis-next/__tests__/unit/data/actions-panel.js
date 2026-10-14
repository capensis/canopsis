import Faker from 'faker';

export const fakeAction = () => ({
  icon: Faker.datatype.string(),
  iconColor: Faker.datatype.string(),
  title: Faker.datatype.string(),
  method: vi.fn(),
});

export const editAction = {
  icon: 'edit',
  iconColor: 'primary',
  title: 'Edit title',
  method: vi.fn(),
};

export const deleteAction = {
  icon: 'delete',
  iconColor: 'secondary',
  title: 'Delete title',
  method: vi.fn(),
};

export const ackAction = {
  icon: 'done',
  iconColor: 'secondary',
  title: 'Ack title',
  method: vi.fn(),
};

export const fastPbehaviorAddAction = {
  icon: 'motion_photos_paused',
  title: 'Fast pbehavior add',
  type: 'fastPbehaviorAdd',
  method: vi.fn(),
};

export const fastPbehaviorRemoveAction = {
  icon: 'play_arrow',
  title: 'Fast pbehavior remove',
  type: 'fastPbehaviorRemove',
  method: vi.fn(),
};
