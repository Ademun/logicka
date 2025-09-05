import {Table, TableBody, TableCell, TableHead, TableHeader, TableRow} from "@/components/ui/table";
import type {TruthTableEntry} from "@/types/truthTable.ts";
import {useState} from "react";
import {Card, CardContent, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {RadioGroup, RadioGroupItem} from "@/components/ui/radio-group.tsx";
import {Label} from "@radix-ui/react-label";

interface TruthTableDisplayProps {
  data: TruthTableEntry[];
}

type HighlightOption = "true" | "false" | "none";

export function TruthTableDisplay({data}: TruthTableDisplayProps) {

  const [highlight, setHighlight] = useState<HighlightOption>("none");

  if (!data || data.length === 0) {
    return <div>Нет данных для отображения</div>;
  }

  // Получаем все уникальные имена переменных из первой записи
  const variableNames = data[0].Variables.map(v => v.Name);

  const getRowClass = (result: boolean) => {
    if (highlight === "none") return "";
    if (highlight === "true" && result) return "bg-green-100";
    if (highlight === "false" && !result) return "bg-red-100";
    return "";
  };

  return (
      <Card>
        <CardHeader>
          <CardTitle>Таблица истинности</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <RadioGroup
              value={highlight}
              onValueChange={(value: HighlightOption) => setHighlight(value)}
              className="flex space-x-4"
          >
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="none" id="r0"/>
              <Label htmlFor="r0">Без подсветки</Label>
            </div>
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="true" id="r2"/>
              <Label htmlFor="r2">Только истина</Label>
            </div>
            <div className="flex items-center space-x-2">
              <RadioGroupItem value="false" id="r3"/>
              <Label htmlFor="r3">Только ложь</Label>
            </div>
          </RadioGroup>

          {/* Таблица с возможностью подсветки */}
          <div className="overflow-x-auto">
            <Table>
              <TableHeader>
                <TableRow>
                  {variableNames.map((name) => (
                      <TableHead key={name} className="font-bold">{name}</TableHead>
                  ))}
                  <TableHead className="font-bold">Результат</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {data.map((entry, index) => (
                    <TableRow
                        key={index}
                        className={getRowClass(entry.Result)}
                    >
                      {entry.Variables.map((variable, varIndex) => (
                          <TableCell key={varIndex}>
                            {variable.Value ? '1' : '0'}
                          </TableCell>
                      ))}
                      <TableCell className="font-medium">
                        {entry.Result ? '1' : '0'}
                      </TableCell>
                    </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>
        </CardContent>
      </Card>

  );
}