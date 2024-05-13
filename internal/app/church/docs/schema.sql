-- ! 1 - 4 USER and DENOMINATION MANAGEMENT
-- Users table 
CREATE TABLE users (
  user_id INT PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password VARCHAR(255) NOT NULL,
  last_denomination_change DATE,
  last_church_change DATE
);

-- User Profiles table
CREATE TABLE user_profiles (
  user_id INT PRIMARY KEY,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  denomination_id INT,
  church_id INT,
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  FOREIGN KEY (denomination_id) REFERENCES denominations(denomination_id),
  FOREIGN KEY (church_id) REFERENCES churches(church_id)
);

-- Denominations table
CREATE TABLE denominations (
  denomination_id INT PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL
);

-- User Denomination Membership table
CREATE TABLE user_denomination_membership (
  user_id INT,
  denomination_id INT,
  join_date DATE,
  PRIMARY KEY (user_id, denomination_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  FOREIGN KEY (denomination_id) REFERENCES denominations(denomination_id)
);

-- Church table
CREATE TABLE churches (
  church_id INT PRIMARY KEY,
  denomination_id INT,
  name VARCHAR(255) NOT NULL,
  members_count int DEFAULT 0,
  FOREIGN KEY (denomination_id) REFERENCES denominations(denomination_id)
);

-- User Church Membership table
CREATE TABLE user_church_membership (
  user_id INT,
  church_id INT,
  join_date DATE,
  PRIMARY KEY (user_id, church_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  FOREIGN KEY (church_id) REFERENCES churches(church_id)
);

-- ! 5 - 6 CHURCH ACCOUNT and NEARBY CHURCHES
-- Church Account and Nearby Churches The churches table defined above will support these features.
-- The location column uses the GEOGRAPHY data type to store spatial data, which can be used to find nearby churches.


-- ! 7 - 11 GROUPS and ORGANIZATION
-- Groups table
CREATE TABLE groups (
  group_id INT PRIMARY KEY,
  denomination_id INT,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  FOREIGN KEY (denomination_id) REFERENCES denominations(denomination_id)
);

-- User Group Membership table
CREATE TABLE user_group_membership (
  user_id INT,
  group_id INT,
  join_date DATE,
  PRIMARY KEY (user_id, group_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  FOREIGN KEY (group_id) REFERENCES groups(group_id)
);

-- Organizations table
CREATE TABLE organizations (
  organization_id INT PRIMARY KEY,
  church_id INT,
  name VARCHAR(255) NOT NULL,
  description TEXT,
  leader_user_id INT,
  FOREIGN KEY (church_id) REFERENCES churches(church_id),
  FOREIGN KEY (leader_user_id) REFERENCES users(user_id)
);

-- User Organization Membership table
CREATE TABLE user_organization_membership (
  user_id INT,
  organization_id INT,
  join_date DATE,
  approved BOOLEAN DEFAULT FALSE,
  PRIMARY KEY (user_id, organization_id),
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  FOREIGN KEY (organization_id) REFERENCES organizations(organization_id)
);


-- ! 12 - 21 FINANCIAL TRX, LEADERSHIP, PLEDGE, NOTIFICATION and WALLET
-- Pledges table
CREATE TABLE pledges (
  pledge_id INT PRIMARY KEY,
  user_id INT,
  church_id INT,
  amount DECIMAL(10, 2),
  pledge_date DATE,
  paid BOOLEAN DEFAULT FALSE,
  FOREIGN KEY (user_id) REFERENCES users(user_id),
  FOREIGN KEY (church_id) REFERENCES churches(church_id)
  );

-- Donations table
CREATE TABLE donations (
    donation_id INT PRIMARY KEY,
    user_id INT,
    church_id INT,
    amount DECIMAL(10, 2),
    donation_type ENUM('tithe', 'offering', 'donation') NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (church_id) REFERENCES churches(church_id)
);

-- Wallets table
CREATE TABLE wallets (
    wallet_id INT PRIMARY KEY,
    user_id INT UNIQUE,
    balance DECIMAL(10, 2) DEFAULT 0.00,
    FOREIGN KEY (user_id) REFERENCES users(user_id)
);

-- Church Bank Accounts table
CREATE TABLE church_bank_accounts (
    church_id INT PRIMARY KEY,
    bank_name VARCHAR(55),
    account_number VARCHAR(255),
    account_name VARCHAR(255),
    FOREIGN KEY (church_id) REFERENCES churches(church_id)
);

-- Transactions table
CREATE TABLE transactions (
    transaction_id INT PRIMARY KEY,
    wallet_id INT,
    amount DECIMAL(10, 2),
    transaction_type ENUM('funding', 'transfer', 'withdrawal') NOT NULL,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (wallet_id) REFERENCES wallets(wallet_id)
);

-- Announcements table
CREATE TABLE announcements (
    announcement_id INT PRIMARY KEY,
    church_id INT,
    content TEXT,
    date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (church_id) REFERENCES churches(church_id)
);

-- Directory table
CREATE TABLE directory (
    directory_id INT PRIMARY KEY,
    church_id INT,
    user_id INT,
    role ENUM('member', 'baptized', 'organization_member') NOT NULL,
    organization_id INT,
    FOREIGN KEY (church_id) REFERENCES churches(church_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (organization_id) REFERENCES organizations(organization_id)
);

-- Mass Bookings table
CREATE TABLE mass_bookings (
    booking_id INT PRIMARY KEY,
    user_id INT,
    church_id INT,
    date DATE,
    last_active TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    FOREIGN KEY (church_id) REFERENCES churches(church_id)
);

-- So for these features 
-- 1. Groups for churches and denominations (eg bible study, SU, etc)
-- 2. Users can join groups irrespective of the denomination, get updates, see and chat with users. 
-- 3. Group forum (Like a group messaging) Users can chat with each other within the group
-- 4. Donate, pay tithe, offering etc to churches (Tithe will be for the church the user is in)
-- 5. Church organizations (CWO, CMO, CYON, Choir etc). Users are approved by admin into the organization. will also have like a group forum for group chatting.