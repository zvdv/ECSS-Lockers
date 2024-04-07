import * as XLSX from "xlsx";
import { db } from "$lib/db";

export async function GET() {
  const data = await db
    .selectFrom("locker")
    .leftJoin("registration", "locker.id", "registration.locker")
    .select([
      "locker.id as locker",
      "registration.name as name",
      "registration.user as email",
    ])
    .orderBy("locker")
    .execute();

  const sheet = XLSX.utils.json_to_sheet(data);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, sheet, "Locker Registration");
  const excelData = await XLSX.writeXLSX(wb, { type: "array" });

  const filename = `locker-registration-${new Date().toISOString()}.xlsx`;
  return new Response(excelData, {
    headers: {
      "Content-Type":
        "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
      "Content-Disposition": `attachment; filename="${filename}"`,
    },
  });
}
