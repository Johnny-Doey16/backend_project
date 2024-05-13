-- Groups table
CREATE TABLE groups (
  id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  denomination_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  FOREIGN KEY (denomination_id) REFERENCES denominations(id)
);
CREATE INDEX idx_groups_denomination_id ON groups(denomination_id);

-- User Group Membership table
CREATE TABLE user_group_membership (
    id INT GENERATED BY DEFAULT AS IDENTITY,
  group_id INT NOT NULL,
  user_id UUID NOT NULL,
  join_date TIMESTAMPTZ DEFAULT now(),
  is_admin BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (id),
  FOREIGN KEY (group_id) REFERENCES groups(id),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);
CREATE INDEX idx_user_group_membership_user_id ON user_group_membership(user_id);
CREATE INDEX idx_user_group_membership_group_id ON user_group_membership(group_id);

-- Organizations table
CREATE TABLE organizations (
  id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  church_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  FOREIGN KEY (church_id) REFERENCES churches(id)
);
CREATE INDEX idx_organizations_church_id ON organizations(church_id);

-- User Organization Membership table
CREATE TABLE user_organization_membership (
    id INT GENERATED BY DEFAULT AS IDENTITY,
  organization_id INT,
  user_id UUID NOT NULL,
  join_date TIMESTAMPTZ DEFAULT now(),
  approved BOOLEAN DEFAULT FALSE,
  is_admin BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (id),
  FOREIGN KEY (organization_id) REFERENCES organizations(id),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);
CREATE INDEX idx_user_organization_membership_user_id ON user_organization_membership(user_id);
CREATE INDEX idx_user_organization_membership_organization_id ON user_organization_membership(organization_id);

-- Group Forums table
CREATE TABLE group_forums (
  id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  group_id INT NOT NULL,
  user_id UUID NOT NULL,
  message TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (group_id) REFERENCES groups(id),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);
CREATE INDEX idx_group_forums_group_id ON group_forums(id);
CREATE INDEX idx_group_forums_user_id ON group_forums(user_id);
CREATE INDEX idx_group_forums_created_at ON group_forums(created_at);

-- Organization Forums table
CREATE TABLE organization_forums (
  id INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  organization_id INT NOT NULL,
  user_id UUID NOT NULL,
  message TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (organization_id) REFERENCES organizations(id),
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX idx_organization_forums_organization_id ON organization_forums(organization_id);
CREATE INDEX idx_organization_forums_user_id ON organization_forums(user_id);
CREATE INDEX idx_organization_forums_created_at ON organization_forums(created_at);