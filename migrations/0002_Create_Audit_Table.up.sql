CREATE TABLE `Audit` (
  `idAduit` INT NOT NULL AUTO_INCREMENT,
  `TITLE` VARCHAR(255) NULL,
  `TEXT` TEXT NULL,
  `DATE` INT NULL,
  PRIMARY KEY (`idAduit`));


  CREATE TABLE `VMS` (
    `idVMS` INT NOT NULL AUTO_INCREMENT,
    `Title` VARCHAR(255) NULL,
    `OS` TEXT NULL,
    `Cores` INT NULL,
    `Memory` INT NULL,
    `Disk` INT NULL,
    `Ipaddress` TEXT NULL,
    PRIMARY KEY (`idVMS`));
