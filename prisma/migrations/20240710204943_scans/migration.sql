-- CreateTable
CREATE TABLE "Scan" (
    "id" SERIAL NOT NULL,
    "userId" INTEGER,
    "url" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Scan_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Scripts" (
    "id" SERIAL NOT NULL,
    "scanId" INTEGER NOT NULL,
    "src" TEXT NOT NULL,
    "status" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Scripts_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "Secrets" (
    "id" SERIAL NOT NULL,
    "secret" TEXT NOT NULL,
    "source" TEXT NOT NULL,
    "line" INTEGER,
    "scriptId" INTEGER NOT NULL,
    "createdAt" TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updatedAt" TIMESTAMP(3) NOT NULL,

    CONSTRAINT "Secrets_pkey" PRIMARY KEY ("id")
);

-- AddForeignKey
ALTER TABLE "Scan" ADD CONSTRAINT "Scan_userId_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE SET NULL ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Scripts" ADD CONSTRAINT "Scripts_scanId_fkey" FOREIGN KEY ("scanId") REFERENCES "Scan"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "Secrets" ADD CONSTRAINT "Secrets_scriptId_fkey" FOREIGN KEY ("scriptId") REFERENCES "Scripts"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
