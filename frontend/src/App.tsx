import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {useEffect, useState} from "react";
import type {TruthTableEntry} from "@/types/truthTable.ts";
import {CalculateTruthTable, ExtractVariables} from "../wailsjs/go/main/App";
import {Card, CardContent, CardDescription, CardHeader, CardTitle} from "@/components/ui/card.tsx";
import {Alert} from "./components/ui/alert";
import {TruthTableDisplay} from "@/components/ui/truthTable.tsx";
import {Label} from "@/components/ui/label.tsx";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue
} from "@/components/ui/select.tsx";

function App() {

  const [logicalExpression, setLogicalExpression] = useState<string>('');
  const [truthTableData, setTruthTableData] = useState<TruthTableEntry[] | null>(null);
  const [error, setError] = useState<string>('');
  const [isLoading, setIsLoading] = useState<boolean>(false);

  const [variables, setVariables] = useState<string[]>([])
  const [variableValues, setVariableValues] = useState<Record<string, boolean | null>>({});

  const evaluateExpression = async () => {
    if (!logicalExpression.trim()) {
      setError('Пожалуйста, введите логическое выражение');
      return;
    }

    setError('');
    setTruthTableData(null);
    setIsLoading(true);

    try {
      const fixedValues: Record<string, boolean> = {};
      Object.entries(variableValues).forEach(([key, value]) => {
        if (value !== null) {
          fixedValues[key] = value;
        }
      });
      const result: TruthTableEntry[] = await CalculateTruthTable(logicalExpression, fixedValues);
      setTruthTableData(result);
    } catch (err: any) {
      setError(err.message || 'Ошибка при вычислении выражения');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (!logicalExpression.trim()) {
      setVariables([]);
      setVariableValues({});
      return;
    }
    setError('')
    extractVars().catch(err => setError(err.message || 'Ошибка при извлечении переменных'))
  }, [logicalExpression]);

  const extractVars = async () => {
    const vars: string[] = await ExtractVariables(logicalExpression)
    if (!vars) {
      setVariables([])
      return
    }
    setVariables(vars)
    // Сохраняем предыдущие значения для существующих переменных
    const newValues: Record<string, boolean | null> = {};
    vars.forEach(v => {
      newValues[v] = variableValues.hasOwnProperty(v) ? variableValues[v] : null;
    });
    setVariableValues(newValues);
  };

  return (
      <div className="container mx-auto p-4 space-y-4 max-w-4xl">
        <Card>
          <CardHeader>
            <CardTitle>Генератор таблиц истинности</CardTitle>
            <CardDescription>
              Введите логическое выражение (например: A && B, !A || C)
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex space-x-2">
              <Input
                  value={logicalExpression}
                  onChange={(e) => setLogicalExpression(e.target.value)}
                  placeholder="Введите логическое выражение..."
                  disabled={isLoading}
                  className="flex-grow"
              />
              <Button
                  onClick={evaluateExpression}
                  disabled={isLoading || !logicalExpression.trim()}
              >
                {isLoading ? "Вычисление..." : "Рассчитать"}
              </Button>
            </div>

            {/* Поля для ввода значений переменных */}
            {(variables && variables.length > 0) && (
                <div className="pt-4 border-t">
                  <h3 className="text-sm font-medium mb-3">Значения переменных (оставьте пустыми для
                    перебора всех значений):</h3>
                  <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
                    {variables.map((variable) => (
                        <div key={variable} className="flex items-center space-x-2">
                          <Label htmlFor={variable} className="w-8">{variable}:</Label>
                          <Select
                              value={variableValues[variable] === null ? "" : String(variableValues[variable])}
                              onValueChange={(value) => {
                                setVariableValues(prev => ({
                                  ...prev,
                                  [variable]: value === "default" ? null : value === "true"
                                }));
                              }}
                          >
                            <SelectTrigger id={variable} className="w-24">
                              <SelectValue placeholder="Все"/>
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="true">True</SelectItem>
                              <SelectItem value="false">False</SelectItem>
                              <SelectItem value="default">Все</SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                    ))}
                  </div>
                </div>
            )}

            {error && (
                <Alert variant="destructive" className="my-4">
                  {error}
                </Alert>
            )}
          </CardContent>
        </Card>

        {/* Компонент для отображения таблицы */}
        {truthTableData && <TruthTableDisplay data={truthTableData}/>}
      </div>
  );
}

export default App
