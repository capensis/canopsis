import Faker from 'faker';

export const fakeAction = () => ({
  icon: Faker.datatype.string(),
  iconColor: Faker.datatype.string(),
  title: Faker.datatype.string(),
  method: jest.fn(),
});

export const editAction = {
  icon: 'edit',
  iconColor: 'primary',
  title: 'Edit title',
  method: jest.fn(),
};

export const deleteAction = {
  icon: 'delete',
  iconColor: 'secondary',
  title: 'Delete title',
  method: jest.fn(),
};

export const ackAction = {
  icon: 'done',
  iconColor: 'secondary',
  title: 'Ack title',
  method: jest.fn(),
};
